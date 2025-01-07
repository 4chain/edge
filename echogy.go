package echogy

import (
	"context"
	"fmt"
	"github.com/echogy-io/echogy/pkg/auth"
	"github.com/echogy-io/echogy/pkg/logger"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"sync"
)

var sessionHub *sync.Map

func init() {
	sessionHub = &sync.Map{}
}

const (
	sshSessionTypeForward            = "tcpip-forward"
	sshSessionTypeCancelForward      = "cancel-tcpip-forward"
	sshRequestTypeDirectTcpip        = "direct-tcpip"
	sshRequestTypeSession            = "session"
	sshAccessIdKey                   = "sshAccessId"
	sshTunnelAddrKey                 = "sshTunnelAddrKey"
	sshDebugServer                   = "sshDebugServer"
	sshRemoteForward                 = "sshRemoteForward"
	clientPublicKeyFingerprintSha256 = "clientPublicKeyFingerprint"
	clientHttpAlias                  = "clientHttpAlias"
	debugPort                        = 4300
)

type remoteForwardSuccess struct {
	BindPort uint32
}

type remoteForwardRequest struct {
	BindAddr string
	BindPort uint32
}

func requestHandler(bindPort uint32) func(ctx ssh.Context, _ *ssh.Server, req *gossh.Request) (bool, []byte) {
	return func(ctx ssh.Context, _ *ssh.Server, req *gossh.Request) (bool, []byte) {
		switch req.Type {
		case sshSessionTypeForward:
			var reqPayload remoteForwardRequest
			if err := gossh.Unmarshal(req.Payload, &reqPayload); err != nil {
				logger.Error("Unmarshal failed", err, map[string]interface{}{
					"module":  "serve",
					"payload": reqPayload,
				})
				return false, []byte{}
			}
			logger.Debug("Unmarshal dispatchRemoteForward request", map[string]interface{}{
				"module":  "serve",
				"payload": reqPayload,
			})
			if reqPayload.BindPort == 0 {
				reqPayload.BindPort = bindPort
			}
			ctx.SetValue(sshRemoteForward, &reqPayload)
			return true, gossh.Marshal(&remoteForwardSuccess{bindPort})

		case sshSessionTypeCancelForward:
			id := ctx.Value(sshAccessIdKey).(string)
			sessionHub.Delete(id)
			return true, nil
		default:
			return false, nil
		}
	}
}

func isRegister(ctx ssh.Context) bool {
	a := ctx.Value(clientHttpAlias)
	return "register" == ctx.User() && nil == a
}

func newSessionServer(sshAddr string, facadeDomain string, sshKey []byte, bindPort uint32, auth auth.Auth) *ssh.Server {
	key, _ := gossh.ParseRawPrivateKey(sshKey)
	signer, _ := gossh.NewSignerFromKey(key)

	reqFunc := requestHandler(bindPort)

	return &ssh.Server{
		//IdleTimeout: 300 * time.Second,
		Version:     "Echogy",
		HostSigners: []ssh.Signer{signer},
		Addr:        sshAddr,
		PtyCallback: func(ctx ssh.Context, pty ssh.Pty) bool {
			return true
		},
		Handler: sessionHandler(facadeDomain),
		PublicKeyHandler: func(ctx ssh.Context, key ssh.PublicKey) bool {
			sha256 := fingerprintSHA256(key)
			if nil != auth {
				alias, found := auth.PubKey(key)
				if found {
					ctx.SetValue(clientPublicKeyFingerprintSha256, sha256)
					ctx.SetValue(clientHttpAlias, alias)
					return true
				}
			}
			if isRegister(ctx) {
				ctx.SetValue(clientPublicKeyFingerprintSha256, sha256)
				return true
			}
			return false
		},
		KeyboardInteractiveHandler: func(ctx ssh.Context, challenger gossh.KeyboardInteractiveChallenge) bool {
			user := ctx.User()
			answers, err := challenger("", "", []string{"Login to Echogy.io\nEnter Password: "}, []bool{false})
			if !isRegister(ctx) && nil != err {
				return false
			}
			if "" != answers[0] {
				alias, found := auth.Password(user, answers[0])
				if found {
					ctx.SetValue(clientHttpAlias, alias)
				}
			}
			return true
		},
		ChannelHandlers: map[string]ssh.ChannelHandler{
			sshRequestTypeDirectTcpip: DirectTCPIPHandler,
			sshRequestTypeSession:     ssh.DefaultSessionHandler,
		},
		RequestHandlers: map[string]ssh.RequestHandler{
			sshSessionTypeForward:       reqFunc,
			sshSessionTypeCancelForward: reqFunc,
		},
	}
}

func sessionHandler(domain string) func(ssh.Session) {
	return func(session ssh.Session) {
		defer func() {
			session.Close()
		}()

		ctx := session.Context()

		var accessId string
		var err error
		var bound bool

		if accessId, bound = ctx.Value(clientHttpAlias).(string); !bound {
			accessId, err = withAddrGenerateAccessId(session.RemoteAddr())
		}

	regenerating:
		if nil != err {
			logger.Error("generating accessId", err, map[string]interface{}{
				"module": "serve",
			})
			session.Write([]byte("generating accessId error"))
			return
		}

		if _, found := sessionHub.Load(accessId); found {
			accessId, err = generateAccessId()
			goto regenerating
		}

		tunnel := fmt.Sprintf("%s.%s", accessId, domain)

		ctx.SetValue(sshTunnelAddrKey, tunnel)

		channel, err := newForwarder(accessId, domain, session)

		if nil != err {
			logger.Error("create dispatchRemoteForward", err, map[string]interface{}{
				"module":     "serve",
				"remoteAddr": session.RemoteAddr().String(),
			})
			return
		}
		ctx.SetValue(sshAccessIdKey, accessId)
		sessionHub.Store(accessId, channel)
		logger.Debug("establishing ssh conn", map[string]interface{}{
			"module":   "serve",
			"accessId": accessId,
		})
		channel.serve() // blocked with loop
		sessionHub.Delete(accessId)
		channel.Close()

		logger.Debug("clean ssh conn", map[string]interface{}{
			"module":   "serve",
			"accessId": accessId,
		})
	}
}

func Serve(ctx context.Context, config *Config, auth auth.Auth) {

	wg := sync.WaitGroup{}

	_, sshPort, err := parseHostAddr(config.SSHAddr)
	if err != nil {
		logger.Fatal("parse net.Addr failed", err, map[string]interface{}{
			"module": "serve",
		})
		return
	}

	server := newSessionServer(config.SSHAddr, config.Domain, []byte(config.PrivateKey), sshPort, auth)

	wg.Add(1)
	go func() {
		wg.Done()
		logger.Warn("started facade server", map[string]interface{}{
			"module":  "serve",
			"address": config.HttpAddr,
		})
		facadeServe(ctx, config.HttpAddr, func(facadeId string, req *hijackHttp) bool {
			if value, found := sessionHub.Load(facadeId); found {
				channel := value.(*forwarder)
				channel.dispatchRemoteForward(req)
				return true
			}
			return false
		})
	}()

	wg.Wait()

	wg.Add(1)
	go func() {
		wg.Done()
		logger.Warn("started ssh server", map[string]interface{}{
			"module":  "serve",
			"address": config.SSHAddr,
		})
		err := server.ListenAndServe()
		logger.Fatal("ssh server shutdown", err, map[string]interface{}{
			"module":  "serve",
			"address": config.SSHAddr,
		})
	}()
	wg.Wait()
	<-ctx.Done()
	server.Shutdown(ctx)
	sessionHub.Range(func(key, value interface{}) bool {
		fwd := value.(*forwarder)
		fwd.Close()
		return true
	})
	logger.WarnN("Echogy shutdown")
}

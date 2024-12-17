
export const Footer = () => {
    return (
        <footer className="mt-auto border-t border-echogy-border dark:border-echogy-border-dark">
            <div className="px-4 py-3">
                <div
                    className="flex items-center space-x-4 text-sm text-echogy-text-secondary dark:text-echogy-text-secondary-dark">
                    <span>{new Date().getFullYear()} <a href="https://echogy.io"
                                                        className="hover:text-echogy-text-primary dark:hover:text-echogy-text-primary-dark transition-colors">Echogy</a></span>
                    <a href="https://echogy.io/privacy" target="_blank" rel="noopener noreferrer"
                       className="hover:text-echogy-text-primary dark:hover:text-echogy-text-primary-dark transition-colors">
                        Privacy
                    </a>
                    <a href="https://echogy.io/terms" target="_blank" rel="noopener noreferrer"
                       className="hover:text-echogy-text-primary dark:hover:text-echogy-text-primary-dark transition-colors">
                        Terms
                    </a>
                </div>
            </div>
        </footer>
    );
};

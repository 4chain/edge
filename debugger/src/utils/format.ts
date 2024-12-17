export const formatTime = (time: number): string => {
    if (time < 1000) {
        return `${time}ms`;
    }
    return `${(time / 1000).toFixed(2)}s`;
};

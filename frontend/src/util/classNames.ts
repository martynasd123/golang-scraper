interface DynamicEntry {
    [key: string]: boolean;
}

type ClassName = DynamicEntry | string;

const classNames = (...args: ClassName[]): string => {
    return args.map(item => {
        if (item == null) {
            return ""
        }
        if (typeof item === "string") {
            return item;
        }
        return Object.keys(item)
            .filter(key => item[key])
            .join(" ");
    }).join(" ");
}

export default classNames;
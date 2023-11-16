type Route = {
    path: string;
    page: any;
};

type Page = {
    title: string;
    page: any;
    path: string;
    isPage: boolean;
};

export default class Docs {
    root: string;
    pages: { [key: string]: Page };
    hierarchy: { [key: string]: any };

    constructor(root: string) {
        this.root = root;
        const imported = import.meta.glob("/docs/**/*.{mdx,yml}", {
            eager: true,
        });
        this.pages = Object.entries(imported).reduce(
            (pages, [fsPath, page]) => {
                const res = { title: "", path: "", isPage: false };
                if (fsPath.endsWith(".yml")) {
                    res.title = page?.label;
                    res.path = this.cleanPath(fsPath);
                } else {
                    res.title = page
                        // @ts-ignore
                        ?.default()
                        ?.props?.children?.find(
                            (child) => child?.type === "h1"
                        )?.props?.children;
                    let path = this.cleanPath(fsPath);
                    if (path.endsWith("/index")) {
                        path = path.replace("/index", "");
                    }
                    res.path = path;
                    res.isPage = true;
                }

                return {
                    ...pages,
                    [res.path]: { ...res, page },
                };
            },
            {}
        );
        this.hierarchy = {};
        this.createHierarchy();
    }

    getRoutes(): Route[] {
        return [
            ...Object.entries(this.pages).map(([path, page]) => {
                return {
                    path: this.cleanPath(path),
                    page: page.page,
                };
            }),
        ];
    }

    getPage(path: string): Page {
        return this.pages[path];
    }

    private cleanPath(filePath: string): string {
        return filePath
            .replace("/docs/", "/")
            .replace("/_category_.yml", "")
            .replace(".mdx", "")
            .split("/")
            .map((segment) => segment.replace(/^\d{2}-/, ""))
            .join("/");
    }

    private createHierarchy() {
        Object.entries(this.pages).forEach(([path]) => {
            const segments = path.split("/").slice(1);
            let group = this.hierarchy;
            segments.forEach((segment) => {
                if (!group[segment]) {
                    group[segment] = {};
                }
                group = group[segment];
            });
            group._path = path;
        });
    }
}

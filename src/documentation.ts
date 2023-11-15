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
                    res.path = fsPath
                        .replace("/docs/", "/")
                        .replace("/_category_.yml", "");
                } else {
                    res.title = page
                        // @ts-ignore
                        ?.default()
                        ?.props?.children?.find(
                            (child) => child?.type === "h1"
                        )?.props?.children;
                    res.path = fsPath
                        .replace("/docs/", "/")
                        .replace(".mdx", "");
                    if (res.path.endsWith("/index")) {
                        res.path = res.path.replace("/index", "");
                    }
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

    getRoutes(): Route[] {
        return [
            ...Object.entries(this.pages).map(([route, page]) => ({
                path: route.replace("/docs/", "/").replace(".mdx", ""),
                page: page.page,
            })),
        ];
    }

    getPage(path: string): Page {
        return this.pages[path];
    }
}

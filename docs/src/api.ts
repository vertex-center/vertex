type Route = {
    path: string;
    tag: string;
    api: any;
};

export default class APIs {
    hierarchy: { [key: string]: any };
    files: { [key: string]: any };
    apis: { [key: string]: any };

    constructor() {
        const imported: any = import.meta.glob("/api/**/*.{yml,yaml}", {
            eager: true,
        });

        this.files = {};
        Object.entries(imported ?? {}).forEach(
            ([fsPath, api]: [string, any]) => {
                this.files[fsPath] = api.default;
            }
        );

        this.apis = {};
        this.hierarchy = {};
        this.populate();
    }

    private populate() {
        const apps: { [key: string]: any } = {};

        Object.entries(this.files).forEach(([fsPath, api]) => {
            const app = fsPath.split("/")[2];
            const version = fsPath.split("/")[3];
            if (!apps[app]) {
                apps[app] = {};
            }
            apps[app][version] = api;
        });

        Object.entries(apps ?? {}).forEach(([app, versions]) => {
            app = "api-" + app;
            Object.entries(versions ?? {}).forEach(([version, api]) => {
                const tags = this.getTags(api);

                tags.forEach((tag: string) => {
                    if (!this.hierarchy[app]) {
                        this.hierarchy[app] = {};
                    }
                    if (!this.hierarchy[app][version]) {
                        this.hierarchy[app][version] = {};
                    }

                    const tagParts = tag.split("/");
                    let path = `/${app}/${version}`;
                    let prev = this.hierarchy[app][version];
                    tagParts.forEach((tagPart: string) => {
                        path += `/${tagPart.replace(/ /g, "-").toLowerCase()}`;
                        if (!prev[tagPart]) {
                            prev[tagPart] = {
                                _path: path,
                                _title: tagPart,
                            };
                        }
                        prev = prev[tagPart];
                    });
                });

                console.log(this.hierarchy);

                tags.forEach((tag: string) => {
                    if (!this.apis[app]) {
                        this.apis[app] = {};
                    }
                    if (!this.apis[app][version]) {
                        this.apis[app][version] = {};
                    }
                    this.apis[app][version][tag] = api;
                });
            });
        });
    }

    private getTags(api: any): Set<string> {
        const tags = new Set<string>();
        const paths = api.paths;
        Object.values(paths ?? {}).forEach((methods: any) => {
            Object.values(methods).forEach((data: any) => {
                if (!data.tags) {
                    return;
                }
                data.tags.forEach((tag: string) => {
                    tags.add(tag);
                });
            });
        });
        return tags;
    }

    getRoutes(): Route[] {
        const routes: Route[] = [];
        Object.entries(this.apis).forEach(([app, versions]) => {
            Object.entries(versions).forEach(
                ([version, tags]: [string, any]) => {
                    routes.push({
                        path: `/${app}/${version}`,
                        tag: "",
                        api: undefined,
                    });
                    Object.entries(tags).forEach(([tag, api]) => {
                        const t = tag.replace(/ /g, "-").toLowerCase();
                        routes.push({
                            path: `/${app}/${version}/${t}`,
                            tag: tag,
                            api: api,
                        });
                    });
                }
            );
        });
        return routes;
    }
}

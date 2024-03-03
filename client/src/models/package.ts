export type Package = {
    name: string;
    description?: string;
    homepage?: string;
    license?: string;
    check?: string;
    install?: { [pm: string]: string };
};

type ProxyRedirect = {
    source: string;
    target: string;
};

type ProxyRedirects = { [uuid: string]: ProxyRedirect };

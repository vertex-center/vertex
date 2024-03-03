declare module "*.sass";

declare global {
    interface Window {
        api_urls: {
            [key: string]: string;
        };
    }
}

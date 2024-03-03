import {
    Bridge,
    Cube,
    Database,
    Pulse,
    ShareNetwork,
    ShieldCheck,
} from "@phosphor-icons/react";

export const useApps = () => {
    return {
        apps: [
            {
                id: "admin",
                name: "Vertex Admin",
                description: "Administer Vertex",
                icon: <ShieldCheck />,
                hidden: false,
            },
            {
                id: "auth",
                name: "Vertex Auth",
                description: "Authentication app for Vertex",
                icon: <ShieldCheck />,
                hidden: true,
            },
            {
                id: "sql",
                name: "Vertex SQL",
                description: "Create and manage SQL databases.",
                icon: <Database />,
                hidden: true,
            },
            {
                id: "tunnels",
                name: "Vertex Tunnels",
                description: "Create and manage tunnels.",
                icon: <Bridge />,
                hidden: true,
            },
            {
                id: "monitoring",
                name: "Vertex Monitoring",
                description: "Create and manage containers.",
                icon: <Pulse />,
                hidden: true,
            },
            {
                id: "containers",
                name: "Vertex Containers",
                description: "Create and manage containers.",
                icon: <Cube />,
                hidden: false,
            },
            {
                id: "reverse-proxy",
                name: "Vertex Reverse Proxy",
                description: "Redirect traffic to your containers.",
                icon: <ShareNetwork />,
                hidden: true,
            },
        ],
    };
};

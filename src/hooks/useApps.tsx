export const useApps = () => {
    return {
        apps: [
            {
                id: "admin",
                name: "Vertex Admin",
                description: "Administer Vertex",
                icon: "admin_panel_settings",
                hidden: false,
            },
            {
                id: "auth",
                name: "Vertex Auth",
                description: "Authentication app for Vertex",
                icon: "admin_panel_settings",
                hidden: true,
            },
            {
                id: "sql",
                name: "Vertex SQL",
                description: "Create and manage SQL databases.",
                icon: "database",
                hidden: true,
            },
            {
                id: "tunnels",
                name: "Vertex Tunnels",
                description: "Create and manage tunnels.",
                icon: "subway",
                hidden: true,
            },
            {
                id: "monitoring",
                name: "Vertex Monitoring",
                description: "Create and manage containers.",
                icon: "monitoring",
                hidden: true,
            },
            {
                id: "containers",
                name: "Vertex Containers",
                description: "Create and manage containers.",
                icon: "deployed_code",
                hidden: false,
            },
            {
                id: "reverse-proxy",
                name: "Vertex Reverse Proxy",
                description: "Redirect traffic to your containers.",
                icon: "router",
                hidden: true,
            },
        ],
    };
};

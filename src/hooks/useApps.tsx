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
                hidden: false,
            },
            {
                id: "tunnels",
                name: "Vertex Tunnels",
                description: "Create and manage tunnels.",
                icon: "subway",
                hidden: false,
            },
            {
                id: "monitoring",
                name: "Vertex Monitoring",
                description: "Create and manage containers.",
                icon: "monitoring",
                hidden: false,
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
                hidden: false,
            },
            {
                id: "devtools-service-editor",
                name: "Vertex Service Editor",
                description: "Create services for publishing.",
                icon: "frame_source",
                category: "devtools",
                hidden: false,
            },
        ],
    };
};

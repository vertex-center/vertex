type VertexApp = {
    id: string;
    icon: string;
    name: string;
    description: string;
    to: string;
};

export const apps: VertexApp[] = [
    {
        id: "vertex-instances",
        to: "/app/vx-instances",
        icon: "storage",
        name: "Vertex Instances",
        description: "Create and manage instances.",
    },
    {
        id: "vertex-monitoring",
        to: "/app/vx-monitoring",
        icon: "monitoring",
        name: "Vertex Monitoring",
        description: "Monitor everything.",
    },
    {
        id: "vertex-tunnels",
        to: "/app/vx-tunnels",
        icon: "subway",
        name: "Vertex Tunnels",
        description: "Create and manage tunnels.",
    },
    {
        id: "vertex-reverse-proxy",
        to: "/app/vx-reverse-proxy",
        icon: "router",
        name: "Vertex Reverse Proxy",
        description: "Redirect traffic to your instances.",
    },
];

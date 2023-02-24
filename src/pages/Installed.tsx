import { useEffect, useState } from "react";
import { getInstalledServices, Service } from "../backend/backend";

type ApplicationProps = {
    service: Service;
};

function Application({ service }: ApplicationProps) {
    return (
        <div className="flex flex-col rounded-md bg-zinc-100 px-4 py-2">
            <h3 className="text-lg font-medium">{service.name}</h3>
        </div>
    );
}

export default function Installed() {
    const [installed, setInstalled] = useState<Service[]>([]);

    useEffect(() => {
        getInstalledServices().then((installed) => setInstalled(installed));
    }, []);

    return (
        <div className="flex flex-col border-separate border-amber-500 gap-4 p-4">
            {installed.map((service) => (
                <Application key={service.id} service={service} />
            ))}
        </div>
    );
}

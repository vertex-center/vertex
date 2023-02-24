import { getAvailableServices, Service } from "../backend/backend";
import { useEffect, useState } from "react";
import Button from "../components/Button";

type ApplicationProps = {
    service: Service;
};

function Application({ service }: ApplicationProps) {
    return (
        <div className="flex flex-col rounded-md bg-zinc-100 px-4 py-2">
            <h3 className="text-lg font-medium">{service.name}</h3>
            <div className="flex flex-row items-center gap-1 text-gray-500">
                <span className="material-symbols-rounded">link</span>
                <a href={`https://${service.repository}`}>
                    {service.repository}
                </a>
            </div>
            <span className="text-sm mt-2">{service.description}</span>
            <div className="flex justify-end">
                <Button rightSymbol="download">Download</Button>
            </div>
        </div>
    );
}

export default function Installed() {
    const [installed, setInstalled] = useState<Service[]>([]);

    useEffect(() => {
        getAvailableServices().then((installed) => setInstalled(installed));
    }, []);

    return (
        <div className="flex flex-col border-separate border-amber-500 gap-4 p-4">
            {installed.map((service) => (
                <Application key={service.id} service={service} />
            ))}
        </div>
    );
}

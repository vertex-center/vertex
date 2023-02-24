import { useEffect, useState } from "react";
import { getInstalledServices, Service } from "../backend/backend";

type ApplicationProps = {
    service: Service;
};

function Application({ service }: ApplicationProps) {
    return (
        <div>
            <h3>{service.name}</h3>
        </div>
    );
}

export default function Installed() {
    const [installed, setInstalled] = useState<Service[]>([]);

    useEffect(() => {
        getInstalledServices().then((installed) => setInstalled(installed));
    }, []);

    return (
        <div>
            {installed.map((service) => (
                <Application key={service.id} service={service} />
            ))}
        </div>
    );
}

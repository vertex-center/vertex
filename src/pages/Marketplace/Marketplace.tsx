import { getAvailableServices, Service } from "../../backend/backend";
import { useEffect, useState } from "react";
import Button from "../../components/Button/Button";
import Card from "../../components/Card/Card";
import { Caption, Title } from "../../components/Text/Text";

import styles from "./Marketplace.module.sass";
import Logo from "../../components/Logo/Logo";
import URL from "../../components/URL/URL";
import Spacer from "../../components/Spacer/Spacer";

type ApplicationProps = {
    service: Service;
};

function Application({ service }: ApplicationProps) {
    return (
        <Card>
            <div className={styles.appHeader}>
                <Logo iconOnly />
                <Title>{service.name}</Title>
                <Spacer />
                <Button rightSymbol="download">Download</Button>
            </div>
            <URL href={`https://${service.repository}`}>
                {service.repository}
            </URL>
            <Caption className={styles.appDescription}>
                {service.description}
            </Caption>
        </Card>
    );
}

export default function Installed() {
    const [installed, setInstalled] = useState<Service[]>([]);

    useEffect(() => {
        getAvailableServices().then((installed) => setInstalled(installed));
    }, []);

    return (
        <div className={styles.cards}>
            {installed.map((service) => (
                <Application key={service.id} service={service} />
            ))}
        </div>
    );
}

import { getAvailableServices, Service } from "../../backend/backend";
import { useEffect, useState } from "react";
import Card from "../../components/Card/Card";
import { Caption, Title } from "../../components/Text/Text";

import styles from "./Marketplace.module.sass";
import Logo from "../../components/Logo/Logo";
import URL from "../../components/URL/URL";
import Spacer from "../../components/Spacer/Spacer";
import classNames from "classnames";
import Button from "../../components/Button/Button";

type ApplicationProps = {
    service: Service;
};

type AppState = "available" | "installing" | "installed";

function Application({ service }: ApplicationProps) {
    const [state, setState] = useState<AppState>("available");

    const onDownloadClick = () => {
        if (state === "available") setState("installing");
        else if (state === "installing") setState("available");
    };

    return (
        <Card>
            <div className={styles.appHeader}>
                <Logo iconOnly />
                <Title>{service.name}</Title>
            </div>
            <URL href={`https://${service.repository}`}>
                {service.repository}
            </URL>
            <Caption className={styles.appDescription}>
                {service.description}
            </Caption>
            <Spacer />
            <div
                className={classNames({
                    [styles.appActions]: true,
                    [styles.downloadBarDownloading]: state === "installing",
                })}
                onClick={onDownloadClick}
            >
                <Button
                    type="transparent"
                    rightSymbol="download"
                    onClick={onDownloadClick}
                >
                    Download
                </Button>
            </div>
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

import {
    getAvailableServices,
    postDownloadService,
    Service,
} from "../../backend/backend";
import { Fragment, useEffect, useState } from "react";
import { Title } from "../../components/Text/Text";

import styles from "./Marketplace.module.sass";
import Button from "../../components/Button/Button";
import Bay from "../../components/Bay/Bay";
import Symbol from "../../components/Symbol/Symbol";
import Select, { Option } from "../../components/Input/Select";
import { Error } from "../../components/Error/Error";
import Loading from "../../components/Loading/Loading";

type DownloadMethod = "marketplace" | "manual";

type StepDownloadProps = {
    onDownload: (service: Service) => void;
};

function StepDownload(props: StepDownloadProps) {
    const { onDownload } = props;

    const [available, setAvailable] = useState<Service[]>([]);

    const [service, setService] = useState<Service>();

    const [method, setMethod] = useState<DownloadMethod>();
    const [error, setError] = useState<string>();

    const [isLoadingMarketplace, setIsLoadingMarketplace] =
        useState<boolean>(false);
    const [isDownloading, setIsDownloading] = useState<boolean>(false);

    useEffect(() => {
        if (method === "marketplace") {
            setIsLoadingMarketplace(true);
            setError(undefined);
            setTimeout(() => {
                getAvailableServices()
                    .then(setAvailable)
                    .catch((err) => {
                        setError(
                            `An error occurred while fetching services from the Marketplace: ${err.message}`
                        );
                        console.error(err);
                    })
                    .finally(() => setIsLoadingMarketplace(false));
            }, 500);
        }
    }, [method]);

    const onServiceChange = (e: any) => {
        let service = available.find((s: Service) => s.id === e.target.value);
        setService(service);
    };

    const download = () => {
        onDownload(service);
        setIsDownloading(true);
    };

    const form = (
        <Fragment>
            <div className={styles.stepTitle}>
                <Symbol name="counter_1" />
                <Title>Download</Title>
            </div>
            <div className={styles.buttons}>
                <Button
                    onClick={() => setMethod("marketplace")}
                    leftSymbol="precision_manufacturing"
                    selectable
                    selected={method === "marketplace"}
                >
                    Marketplace
                </Button>
                <Button
                    onClick={() => setMethod("manual")}
                    leftSymbol="hand_gesture"
                    selectable
                    selected={method === "manual"}
                >
                    Manual
                </Button>
            </div>
            {method === "marketplace" && !isLoadingMarketplace && !error && (
                <Select label="Service" onChange={onServiceChange}>
                    <Option />
                    {available.map((service) => (
                        <Option key={service.id} value={service.id}>
                            {service.name}
                        </Option>
                    ))}
                </Select>
            )}
            {method === "marketplace" && isLoadingMarketplace && <Loading />}
            <Button
                primary
                large
                rightSymbol="download"
                disabled={!service}
                onClick={download}
            >
                Download
            </Button>
        </Fragment>
    );

    return (
        <div className={styles.step}>
            {!isDownloading && form}
            <Error error={error} />
        </div>
    );
}

export default function Installed() {
    const [isDownloading, setIsDownloading] = useState<boolean>(false);
    const [service, setService] = useState<Service>();

    const download = (service: Service) => {
        setIsDownloading(true);
        setService(service);
        postDownloadService(service).then(() => {
            setIsDownloading(false);
        });
    };

    return (
        <div className={styles.marketplace}>
            <div className={styles.content}>
                <div className={styles.server}>
                    <Bay
                        name={service?.name ?? "Empty server"}
                        status={isDownloading ? "downloading" : "off"}
                    />
                </div>
                <StepDownload onDownload={download} />
            </div>
        </div>
    );
}

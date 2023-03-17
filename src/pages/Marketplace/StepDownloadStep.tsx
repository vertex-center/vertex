import { getAvailableServices, Service } from "../../backend/backend";
import { Fragment, useEffect, useState } from "react";
import styles from "./Marketplace.module.sass";
import { Title } from "../../components/Text/Text";
import Select, { Option } from "../../components/Input/Select";
import Loading from "../../components/Loading/Loading";
import Button from "../../components/Button/Button";
import { Error } from "../../components/Error/Error";
import { DownloadMethod } from "./Marketplace";

type StepDownloadProps = {
    method: DownloadMethod;
    onDownload: (service: Service) => void;
};

export default function StepDownload(props: StepDownloadProps) {
    const { method, onDownload } = props;

    const [available, setAvailable] = useState<Service[]>([]);

    const [service, setService] = useState<Service>();

    const [error, setError] = useState<string>();

    const [isLoadingMarketplace, setIsLoadingMarketplace] =
        useState<boolean>(false);
    const [isDownloading, setIsDownloading] = useState<boolean>(false);

    useEffect(() => {
        if (method === "marketplace") {
            setIsLoadingMarketplace(true);
            setError(undefined);
            getAvailableServices()
                .then(setAvailable)
                .catch((err) => {
                    setError(
                        `An error occurred while fetching services from the Marketplace: ${err.message}`
                    );
                    console.error(err);
                })
                .finally(() => setIsLoadingMarketplace(false));
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
                <Title>Download</Title>
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
            {method === "localstorage" && <></>}
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

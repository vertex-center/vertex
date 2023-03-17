import { getAvailableServices, Service } from "../../backend/backend";
import { Fragment, useEffect, useState } from "react";
import styles from "./Marketplace.module.sass";
import { Title } from "../../components/Text/Text";
import Select, { Option } from "../../components/Input/Select";
import Loading from "../../components/Loading/Loading";
import Button from "../../components/Button/Button";
import { Error } from "../../components/Error/Error";
import { DownloadMethod } from "./Marketplace";

type StepDownloadMarketplaceProps = {
    onNextStep: () => void;
    service: Service;
    onServiceChange: (service: Service) => void;
};

function StepDownloadMarketplace(props: StepDownloadMarketplaceProps) {
    const { service, onNextStep } = props;

    const [available, setAvailable] = useState<Service[]>([]);

    const [error, setError] = useState<string>();

    const [isLoadingMarketplace, setIsLoadingMarketplace] =
        useState<boolean>(false);

    useEffect(() => {
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
    }, []);

    const onServiceChange = (e: any) => {
        let service = available.find((s: Service) => s.id === e.target.value);
        props.onServiceChange(service);
    };

    return (
        <Fragment>
            {!isLoadingMarketplace && !error && (
                <Select label="Service" onChange={onServiceChange}>
                    <Option />
                    {available.map((service) => (
                        <Option key={service.id} value={service.id}>
                            {service.name}
                        </Option>
                    ))}
                </Select>
            )}
            {isLoadingMarketplace && <Loading />}
            <Button
                primary
                large
                rightSymbol="download"
                disabled={!service}
                onClick={onNextStep}
            >
                Download
            </Button>
            <Error error={error} />
        </Fragment>
    );
}

type StepDownloadLocalStorageProps = {
    onNextStep: () => void;
    service: Service;
    onServiceChange: (service: Service) => void;
};

function StepDownloadLocalStorage(props: StepDownloadLocalStorageProps) {
    return null;
}

type StepDownloadProps = {
    method: DownloadMethod;
    service: Service;
    onServiceChange: (service: Service) => void;
    onNextStep: () => void;
};

export default function StepDownload(props: StepDownloadProps) {
    const { method, service, onServiceChange, onNextStep } = props;

    return (
        <div className={styles.step}>
            <div className={styles.stepTitle}>
                <Title>Download</Title>
            </div>
            {method === "marketplace" && (
                <StepDownloadMarketplace
                    onNextStep={onNextStep}
                    service={service}
                    onServiceChange={onServiceChange}
                />
            )}
            {method === "localstorage" && (
                <StepDownloadLocalStorage
                    onNextStep={onNextStep}
                    service={service}
                    onServiceChange={onServiceChange}
                />
            )}
        </div>
    );
}

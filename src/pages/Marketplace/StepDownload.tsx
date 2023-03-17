import { getAvailableServices, Service } from "../../backend/backend";
import { Fragment, useEffect, useState } from "react";
import styles from "./Marketplace.module.sass";
import { Title } from "../../components/Text/Text";
import Select, { Option } from "../../components/Input/Select";
import Loading from "../../components/Loading/Loading";
import Button from "../../components/Button/Button";
import { Error } from "../../components/Error/Error";
import { DownloadMethod } from "./Marketplace";
import Input from "../../components/Input/Input";

type StepDownloadMarketplaceProps = {
    onNextStep: () => void;
    repository: string;
    onRepositoryChange: (repository: string) => void;
};

function StepDownloadMarketplace(props: StepDownloadMarketplaceProps) {
    const { repository, onNextStep } = props;

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
        props.onRepositoryChange(e.target.value);
    };

    return (
        <Fragment>
            {!isLoadingMarketplace && !error && (
                <Select label="Service" onChange={onServiceChange}>
                    <Option />
                    {available.map((service) => (
                        <Option key={service.id} value={service.repository}>
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
                disabled={!repository}
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
    repository: string;
    onRepositoryChange: (repository: string) => void;
};

function StepDownloadLocalStorage(props: StepDownloadLocalStorageProps) {
    const { onNextStep, repository, onRepositoryChange } = props;

    const [path, setPath] = useState();

    const onPathChange = (e: any) => {
        setPath(e.target.value);
        onRepositoryChange(`localstorage:${e.target.value}`);
    };

    return (
        <Fragment>
            <Input
                value={path}
                onChange={onPathChange}
                label="Service path"
                description="Absolute path on your local machine"
            />
            <Button
                primary
                large
                rightSymbol="link"
                disabled={!repository}
                onClick={onNextStep}
            >
                Link to Vertex
            </Button>
        </Fragment>
    );
}

type StepDownloadProps = {
    method: DownloadMethod;
    repository: string;
    onRepositoryChange: (repository: string) => void;
    onNextStep: () => void;
};

export default function StepDownload(props: StepDownloadProps) {
    const { method, repository, onRepositoryChange, onNextStep } = props;

    return (
        <div className={styles.step}>
            <div className={styles.stepTitle}>
                <Title>Download</Title>
            </div>
            {method === "marketplace" && (
                <StepDownloadMarketplace
                    onNextStep={onNextStep}
                    repository={repository}
                    onRepositoryChange={onRepositoryChange}
                />
            )}
            {method === "localstorage" && (
                <StepDownloadLocalStorage
                    onNextStep={onNextStep}
                    repository={repository}
                    onRepositoryChange={onRepositoryChange}
                />
            )}
        </div>
    );
}

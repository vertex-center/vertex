import { getAvailableServices, Service } from "../../backend/backend";
import { Fragment, useEffect, useState } from "react";
import styles from "./Marketplace.module.sass";
import { Text, Title } from "../../components/Text/Text";
import Select, { Option } from "../../components/Input/Select";
import Loading from "../../components/Loading/Loading";
import Button from "../../components/Button/Button";
import { Error } from "../../components/Error/Error";
import { DownloadMethod } from "./Marketplace";
import Input from "../../components/Input/Input";
import { SegmentedButtons } from "../../components/SegmentedButton";
import { Horizontal } from "../../components/Layouts/Layouts";
import Spacer from "../../components/Spacer/Spacer";

type StepDownloadMarketplaceProps = {
    onNextStep: () => void;
    repository: string;
    useDocker?: boolean;
    useReleases?: boolean;
    onRepositoryChange: (repository: string) => void;
    onUseDockerChange: (useDocker: boolean) => void;
    onUseReleasesChange: (useReleases: boolean) => void;
};

function StepDownloadMarketplace(props: StepDownloadMarketplaceProps) {
    const {
        onNextStep,
        repository,
        useDocker,
        useReleases,
        onUseDockerChange,
        onUseReleasesChange,
    } = props;

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
        console.log(e.target.value);
        if (e.target.value === "UNDEFINED")
            return props.onRepositoryChange(undefined);

        props.onRepositoryChange(`marketplace:${e.target.value}`);
    };

    return (
        <Fragment>
            {!isLoadingMarketplace && !error && (
                <Select label="Service" onChange={onServiceChange}>
                    <Option value="UNDEFINED" />
                    {available.map((service) => (
                        <Option key={service.id} value={service.repository}>
                            {service.name}
                        </Option>
                    ))}
                </Select>
            )}
            <Horizontal alignItems="center">
                <Text>Use Docker?</Text>
                <Spacer />
                <SegmentedButtons
                    value={useDocker}
                    onChange={(v) => onUseDockerChange(v)}
                    items={[
                        {
                            label: "Yes",
                            value: true,
                            rightSymbol: "check",
                        },
                        {
                            label: "No",
                            value: false,
                            rightSymbol: "close",
                        },
                    ]}
                />
            </Horizontal>
            {useDocker === false && (
                <Horizontal alignItems="center">
                    <Text>Download the precompiled release?</Text>
                    <Spacer />
                    <SegmentedButtons
                        value={useReleases}
                        onChange={(v) => onUseReleasesChange(v)}
                        items={[
                            {
                                label: "Yes",
                                value: true,
                                rightSymbol: "check",
                            },
                            {
                                label: "No",
                                value: false,
                                rightSymbol: "close",
                            },
                        ]}
                    />
                </Horizontal>
            )}
            {isLoadingMarketplace && <Loading />}
            <Button
                primary
                large
                rightSymbol="download"
                disabled={
                    !repository ||
                    useDocker === undefined ||
                    (useReleases === undefined && useDocker === false)
                }
                onClick={onNextStep}
            >
                Download
            </Button>
            <Error error={error} />
        </Fragment>
    );
}

type StepDownloadGitHubProps = {
    onNextStep: () => void;
    repository: string;
    onRepositoryChange: (repository: string) => void;
};

function StepDownloadGitHub(props: StepDownloadGitHubProps) {
    const { repository, onNextStep } = props;

    const onLinkChange = (e: any) => {
        props.onRepositoryChange(`git:${e.target.value}`);
    };

    return (
        <Fragment>
            <Input
                label="Repository"
                placeholder="https://github.com/..."
                onChange={onLinkChange}
            />
            <Button
                primary
                large
                rightSymbol="download"
                disabled={!repository}
                onClick={onNextStep}
            >
                Download
            </Button>
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
    useDocker?: boolean;
    useReleases?: boolean;
    onRepositoryChange: (repository: string) => void;
    onUseDockerChange: (useDocker: boolean) => void;
    onUseReleasesChange: (useReleases: boolean) => void;
    onNextStep: () => void;
};

export default function StepDownload(props: StepDownloadProps) {
    const {
        method,
        repository,
        useDocker,
        useReleases,
        onRepositoryChange,
        onUseDockerChange,
        onUseReleasesChange,
        onNextStep,
    } = props;

    return (
        <div className={styles.step}>
            <div className={styles.stepTitle}>
                <Title>Download</Title>
            </div>
            {method === "marketplace" && (
                <StepDownloadMarketplace
                    onNextStep={onNextStep}
                    repository={repository}
                    useDocker={useDocker}
                    useReleases={useReleases}
                    onRepositoryChange={onRepositoryChange}
                    onUseDockerChange={onUseDockerChange}
                    onUseReleasesChange={onUseReleasesChange}
                />
            )}
            {method === "git" && (
                <StepDownloadGitHub
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

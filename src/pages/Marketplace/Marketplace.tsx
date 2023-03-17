import {
    downloadService,
    EnvVariable,
    getAvailableServices,
    Service,
} from "../../backend/backend";
import { Fragment, useEffect, useState } from "react";
import { Caption, Title } from "../../components/Text/Text";

import styles from "./Marketplace.module.sass";
import Button from "../../components/Button/Button";
import Bay from "../../components/Bay/Bay";
import Symbol from "../../components/Symbol/Symbol";
import Select, { Option } from "../../components/Input/Select";
import { Error } from "../../components/Error/Error";
import Loading from "../../components/Loading/Loading";
import Input from "../../components/Input/Input";
import { Vertical } from "../../components/Layouts/Layouts";
import PortInput from "../../components/Input/PortInput";

type DownloadMethod = "marketplace" | "localstorage";

type StepSelectMethodProps = {
    method: DownloadMethod;
    onMethodChange: (method: DownloadMethod) => void;
    onNextStep: () => void;
};

function StepSelectMethod(props: StepSelectMethodProps) {
    const { method, onMethodChange, onNextStep } = props;

    return (
        <div className={styles.step}>
            <div className={styles.stepTitle}>
                <Title>Installation method</Title>
            </div>
            <div className={styles.buttons}>
                <Button
                    className={styles.button}
                    onClick={() => onMethodChange("marketplace")}
                    leftSymbol="precision_manufacturing"
                    selectable
                    selected={method === "marketplace"}
                >
                    <div className={styles.buttonContent}>
                        <div>Marketplace</div>
                        <div className={styles.buttonDescription}>
                            Download services from our online and certified
                            repository.
                        </div>
                    </div>
                </Button>
                <Button
                    className={styles.button}
                    onClick={() => onMethodChange("localstorage")}
                    leftSymbol="storage"
                    selectable
                    selected={method === "localstorage"}
                >
                    <div className={styles.buttonContent}>
                        <div>Local storage</div>
                        <div className={styles.buttonDescription}>
                            Point to vertex the path of an existing service on
                            your computer. Vertex will keep it installed there.
                        </div>
                    </div>
                </Button>
            </div>
            <Button
                primary
                large
                disabled={method === undefined}
                rightSymbol="navigate_next"
                onClick={onNextStep}
            >
                Next
            </Button>
        </div>
    );
}

type StepDownloadProps = {
    method: DownloadMethod;
    onDownload: (service: Service) => void;
};

function StepDownload(props: StepDownloadProps) {
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

type VariableInputProps = {
    env: EnvVariable;
    value: any;
    onChange: (value: any) => void;
};

function VariableInput(props: VariableInputProps) {
    const { env, value, onChange } = props;

    const inputProps = {
        value,
        label: env.display_name,
        name: env.name,
        onChange: (e) => onChange(e.target.value),
    };

    let input;
    if (env.type === "port") {
        input = <PortInput {...inputProps} />;
    } else {
        input = <Input {...inputProps} />;
    }

    return (
        <Vertical gap={6}>
            {input}
            <Caption className={styles.inputDescription}>
                {env.description}
            </Caption>
        </Vertical>
    );
}

type StepConfigureProps = {
    service: Service;
};

function StepConfigure(props: StepConfigureProps) {
    const { service } = props;

    const [env, setEnv] = useState<any[]>();

    useEffect(() => {
        setEnv(
            service.environment.map((e) => ({
                env: e,
                value: e.default ?? "",
            }))
        );
    }, [service.environment]);

    const onChange = (i: number, value: any) => {
        setEnv((prev) =>
            prev.map((el, index) => {
                if (index !== i) return el;
                return { ...el, value };
            })
        );
    };

    return (
        <div className={styles.step}>
            <div className={styles.stepTitle}>
                <Symbol name="counter_2" />
                <Title>Configure</Title>
            </div>
            <Vertical gap={30}>
                {env?.map((e, i) => (
                    <VariableInput
                        key={i}
                        env={e.env}
                        value={e.value}
                        onChange={(v: any) => onChange(i, v)}
                    />
                ))}
            </Vertical>
        </div>
    );
}

type Step = "select-method" | "download" | "downloading" | "configure";

export default function Installed() {
    const [step, setStep] = useState<Step>("select-method");

    const [service, setService] = useState<Service>();
    const [method, setMethod] = useState<DownloadMethod>();

    const [error, setError] = useState<string>();

    const download = (service: Service) => {
        setStep("downloading");
        downloadService(service)
            .then((data: any) => {
                console.log(data.instance);
                setStep("configure");
                setService(data.instance);
            })
            .catch((error) => {
                console.log(error);
                setError(`${error.message}: ${error.response.data.message}`);
            });
    };

    let status;
    switch (step) {
        case "download":
            status = "off";
            break;
        case "downloading":
            status = "downloading";
            break;
        case "configure":
            status = "waiting";
            break;
    }

    return (
        <div className={styles.marketplace}>
            <div className={styles.content}>
                <div className={styles.server}>
                    {step === "downloading" && !error && (
                        <Fragment>
                            <div className={styles.cloud}>
                                <Symbol name="cloud" />
                            </div>
                            <div className={styles.cable}></div>
                        </Fragment>
                    )}
                    <Bay
                        name={service?.name ?? "Empty server"}
                        status={status ?? "off"}
                    />
                </div>
                {step === "select-method" && (
                    <StepSelectMethod
                        method={method}
                        onMethodChange={(m) => setMethod(m)}
                        onNextStep={() => setStep("download")}
                    />
                )}
                {step === "download" && (
                    <StepDownload method={method} onDownload={download} />
                )}
                {step === "configure" && <StepConfigure service={service} />}
                <Error error={error} />
            </div>
        </div>
    );
}

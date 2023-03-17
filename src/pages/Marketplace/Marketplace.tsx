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

type Step = "download" | "downloading" | "configure";

export default function Installed() {
    const [step, setStep] = useState<Step>("download");

    const [service, setService] = useState<Service>();

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
                        status={status}
                    />
                </div>
                {step === "download" && <StepDownload onDownload={download} />}
                {step === "configure" && <StepConfigure service={service} />}
                <Error error={error} />
            </div>
        </div>
    );
}

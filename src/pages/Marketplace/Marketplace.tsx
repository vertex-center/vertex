import { downloadService, Service } from "../../backend/backend";
import { Fragment, useState } from "react";

import styles from "./Marketplace.module.sass";
import Bay from "../../components/Bay/Bay";
import Symbol from "../../components/Symbol/Symbol";
import { Error } from "../../components/Error/Error";
import StepSelectMethod from "./StepSelectMethod";
import StepDownload from "./StepDownload";
import StepConfigure from "./StepConfigure";

export type DownloadMethod = "marketplace" | "localstorage";

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
                    <StepDownload
                        method={method}
                        service={service}
                        onServiceChange={(s) => setService(s)}
                        onNextStep={() => download(service)}
                    />
                )}
                {step === "configure" && <StepConfigure service={service} />}
                <Error error={error} />
            </div>
        </div>
    );
}

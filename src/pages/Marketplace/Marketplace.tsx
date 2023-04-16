import { downloadService, Instance } from "../../backend/backend";
import { Fragment, useState } from "react";

import styles from "./Marketplace.module.sass";
import Bay from "../../components/Bay/Bay";
import Symbol from "../../components/Symbol/Symbol";
import { Error } from "../../components/Error/Error";
import StepSelectMethod from "./StepSelectMethod";
import StepDownload from "./StepDownload";
import StepConfigure from "./StepConfigure";
import { useNavigate } from "react-router-dom";
import Header from "../../components/Header/Header";
import { SiGithub } from "@icons-pack/react-simple-icons";

export type DownloadMethod = "marketplace" | "git" | "localstorage";

type Step = "select-method" | "download" | "downloading" | "configure";

export default function Installed() {
    const navigate = useNavigate();

    const [step, setStep] = useState<Step>("select-method");

    const [repository, setRepository] = useState<string>();
    const [useDocker, setUseDocker] = useState<boolean>();
    const [useReleases, setUseReleases] = useState<boolean>();

    const [instance, setInstance] = useState<Instance>();
    const [method, setMethod] = useState<DownloadMethod>();

    const [error, setError] = useState<string>();

    const download = (
        repository: string,
        useDocker?: boolean,
        useReleases?: boolean
    ) => {
        setStep("downloading");
        downloadService(repository, useDocker, useReleases)
            .then((data: any) => {
                console.log(data.instance);
                setStep("configure");
                setInstance(data.instance);
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
        <div>
            <Header />
            <div className={styles.marketplace}>
                <div className={styles.content}>
                    <div className={styles.server}>
                        <Bay
                            instances={[
                                {
                                    name: instance?.name ?? "Empty server",
                                    status: status ?? "off",
                                },
                            ]}
                        />
                        {step === "downloading" && !error && (
                            <Fragment>
                                <div className={styles.cable}></div>
                                <div className={styles.cloud}>
                                    {method === "localstorage" ? (
                                        <Symbol name="storage" />
                                    ) : (
                                        <SiGithub />
                                    )}
                                </div>
                            </Fragment>
                        )}
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
                            repository={repository}
                            useDocker={useDocker}
                            useReleases={useReleases}
                            onRepositoryChange={(v) => setRepository(v)}
                            onUseDockerChange={(v) => setUseDocker(v)}
                            onUseReleasesChange={(v) => setUseReleases(v)}
                            onNextStep={() =>
                                download(repository, useDocker, useReleases)
                            }
                        />
                    )}
                    {step === "configure" && (
                        <StepConfigure
                            onNextStep={() =>
                                navigate(`/bay/${instance.uuid}/`)
                            }
                            instance={instance}
                        />
                    )}
                    <Error error={error} />
                </div>
            </div>
        </div>
    );
}

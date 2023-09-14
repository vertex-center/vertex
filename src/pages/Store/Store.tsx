import { useFetch } from "../../hooks/useFetch";
import { BigTitle, Text, Title } from "../../components/Text/Text";
import { Service as ServiceModel } from "../../models/service";
import { Fragment, useState } from "react";

import styles from "./Store.module.sass";
import Service from "../../components/Service/Service";
import { Error } from "../../components/Error/Error";
import { Horizontal, Vertical } from "../../components/Layouts/Layouts";
import Spacer from "../../components/Spacer/Spacer";
import { SegmentedButtons } from "../../components/SegmentedButton";
import Button from "../../components/Button/Button";
import Popup from "../../components/Popup/Popup";
import { InstallMethod, Instances } from "../../models/instance";
import { api } from "../../backend/backend";

type Downloading = {
    service: ServiceModel;
};

// type ImportMethod = "git" | "localstorage";

export default function Store() {
    const { data: services } = useFetch<ServiceModel[]>(
        api.services.available.get
    );
    const { data: instances, reload: reloadInstances } = useFetch<Instances>(
        api.instances.get
    );

    const [installMethod, setInstallMethod] = useState<InstallMethod>();

    const [showInstallPopup, setShowInstallPopup] = useState<boolean>(false);
    const [showImportPopup, setShowImportPopup] = useState<boolean>(false);

    // const [repository, setRepository] = useState();
    // const [importMethod, setImportMethod] = useState<ImportMethod>("git");
    // const [importing, setImporting] = useState(false);

    const [selectedService, setSelectedService] = useState<ServiceModel>();

    const [error, setError] = useState<string>();
    const [popupError, setPopupError] = useState<string>();

    const [downloading, setDownloading] = useState<Downloading[]>([]);

    const openInstallPopup = (service: ServiceModel) => {
        setSelectedService(service);
        setShowInstallPopup(true);
        setInstallMethod("script");
        setPopupError(undefined);

        const { script, release, docker } = service.methods;

        if (script) {
            setInstallMethod("script");
        } else if (release) {
            setInstallMethod("release");
        } else if (docker) {
            setInstallMethod("docker");
        } else {
            setPopupError("This service doesn't have installation method.");
        }
    };

    const closeInstallPopup = () => {
        setSelectedService(undefined);
        setShowInstallPopup(false);
    };

    const install = () => {
        // UUID to remove from the download queue after installation
        const service = selectedService;

        setDownloading((prev) => [...prev, { service }]);
        setShowInstallPopup(false);

        api.services
            .install({
                method: installMethod,
                service_id: service.id,
            })
            .catch((error) => {
                setError(error.message);
            })
            .finally(() => {
                setDownloading(
                    downloading.filter(({ service: s }) => s.id !== service.id)
                );
                reloadInstances().catch(console.error);
            });
    };

    // const installFromElsewhere = () => {
    //     setImporting(true);
    //     installService(`${importMethod}:${repository}`)
    //         .catch((error) => {
    //             setError(error.message);
    //         })
    //         .finally(() => {
    //             setShowImportPopup(false);
    //             setImporting(false);
    //         });
    // };

    // const onRepoChange = (e: any) => {
    //     setRepository(e.target.value);
    // };

    const installPopup = (
        <Popup
            show={showInstallPopup}
            onDismiss={() => setShowInstallPopup(false)}
        >
            <Title>Download {selectedService?.name}</Title>
            <Error error={popupError} />
            {!popupError && (
                <Horizontal alignItems="center" gap={12}>
                    <Text>Installation method</Text>
                    <Spacer />
                    <SegmentedButtons
                        value={installMethod}
                        onChange={(v) => setInstallMethod(v)}
                        items={[
                            {
                                label: "Script",
                                value: "script",
                                rightSymbol: "description",
                                disabled: !selectedService?.methods?.script,
                            },
                            {
                                label: "Release",
                                value: "release",
                                rightSymbol: "package",
                                disabled: !selectedService?.methods?.release,
                            },
                            {
                                label: "Docker",
                                value: "docker",
                                rightSymbol: "deployed_code",
                                disabled: !selectedService?.methods?.docker,
                            },
                        ]}
                    />
                </Horizontal>
            )}
            <Horizontal gap={8}>
                <Spacer />
                <Button onClick={closeInstallPopup}>Cancel</Button>
                <Button
                    onClick={install}
                    primary
                    rightSymbol="download"
                    disabled={popupError !== undefined}
                >
                    Download
                </Button>
            </Horizontal>
        </Popup>
    );

    // const importPopup = (
    //     <Popup
    //         show={showImportPopup}
    //         onDismiss={() => setShowImportPopup(false)}
    //     >
    //         <Title>Import from elsewhere</Title>
    //         <SegmentedButtons
    //             disabled={importing}
    //             value={importMethod}
    //             onChange={(v) => setImportMethod(v)}
    //             items={[
    //                 {
    //                     label: "Git remote",
    //                     value: "git",
    //                     rightSymbol: "cloud_download",
    //                 },
    //                 {
    //                     label: "Local storage",
    //                     value: "localstorage",
    //                     rightSymbol: "storage",
    //                 },
    //             ]}
    //         />
    //         {importMethod === "git" && (
    //             <Input
    //                 disabled={importing}
    //                 value={repository}
    //                 onChange={onRepoChange}
    //                 label="Repository"
    //                 placeholder="github.com/user/repo"
    //                 description="All Git remotes are compatible."
    //             />
    //         )}
    //         {importMethod === "localstorage" && (
    //             <Input
    //                 disabled={importing}
    //                 value={repository}
    //                 onChange={onRepoChange}
    //                 label="Service path"
    //                 description="Absolute path on your local machine where the server is running."
    //             />
    //         )}
    //         {importing && <Progress infinite />}
    //         <Horizontal gap={12}>
    //             <Spacer />
    //             <Button
    //                 disabled={importing}
    //                 onClick={() => setShowImportPopup(false)}
    //             >
    //                 Cancel
    //             </Button>
    //             <Button
    //                 disabled={importing}
    //                 onClick={() => installFromElsewhere()}
    //                 primary
    //                 rightSymbol={importMethod === "git" ? "download" : "link"}
    //             >
    //                 {importMethod === "git" ? "Download" : "Link"}
    //             </Button>
    //         </Horizontal>
    //     </Popup>
    // );

    return (
        <Fragment>
            <div className={styles.page}>
                <Horizontal
                    className={styles.title}
                    gap={10}
                    alignItems="center"
                >
                    <BigTitle>Marketplace</BigTitle>
                </Horizontal>
                {<Error error={error} />}
                <Vertical className={styles.content}>
                    {services?.map((service) => (
                        <Service
                            key={service.id}
                            service={service}
                            onInstall={() => openInstallPopup(service)}
                            downloading={downloading.some(
                                ({ service: s }) => s.id === service.id
                            )}
                            installedCount={
                                instances === undefined
                                    ? undefined
                                    : Object.values(instances)?.filter(
                                          ({ id }) => id === service.id
                                      )?.length
                            }
                        />
                    ))}
                </Vertical>
            </div>
            {installPopup}
            {/*{importPopup}*/}
        </Fragment>
    );
}

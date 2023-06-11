import { useFetch } from "../../hooks/useFetch";
import {
    downloadService,
    getAvailableServices,
    getInstances,
} from "../../backend/backend";
import { Text, Title } from "../../components/Text/Text";
import { Service as ServiceModel } from "../../models/service";
import Header from "../../components/Header/Header";
import { Fragment, useState } from "react";

import styles from "./Store.module.sass";
import Service from "../../components/Service/Service";
import { Error } from "../../components/Error/Error";
import { Horizontal, Vertical } from "../../components/Layouts/Layouts";
import Spacer from "../../components/Spacer/Spacer";
import { SegmentedButtons } from "../../components/SegmentedButton";
import Button from "../../components/Button/Button";
import Popup from "../../components/Popup/Popup";
import { Instances } from "../../models/instance";
import Input from "../../components/Input/Input";
import Progress from "../../components/Progress";

type Downloading = {
    service: ServiceModel;
};

type ImportMethod = "git" | "localstorage";

export default function Store() {
    const { data: services } = useFetch<ServiceModel[]>(getAvailableServices);
    const { data: instances, reload: reloadInstances } =
        useFetch<Instances>(getInstances);

    const [useDocker, setUseDocker] = useState<boolean>(false);
    const [useReleases, setUseReleases] = useState<boolean>(false);

    const [showInstallPopup, setShowInstallPopup] = useState<boolean>(false);
    const [showImportPopup, setShowImportPopup] = useState<boolean>(false);

    const [repository, setRepository] = useState();
    const [importMethod, setImportMethod] = useState<ImportMethod>("git");
    const [importing, setImporting] = useState(false);

    const [selectedService, setSelectedService] = useState<ServiceModel>();

    const [error, setError] = useState();

    const [downloading, setDownloading] = useState<Downloading[]>([]);

    const openInstallPopup = (service: ServiceModel) => {
        setSelectedService(service);
        setShowInstallPopup(true);
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

        downloadService(
            `marketplace:${service.repository}`,
            useDocker,
            useReleases
        )
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

    const installFromElsewhere = () => {
        setImporting(true);
        downloadService(`${importMethod}:${repository}`)
            .catch((error) => {
                setError(error.message);
            })
            .finally(() => {
                setShowImportPopup(false);
                setImporting(false);
            });
    };

    const onRepoChange = (e: any) => {
        setRepository(e.target.value);
    };

    const installPopup = (
        <Popup
            show={showInstallPopup}
            onDismiss={() => setShowInstallPopup(false)}
        >
            <Title>Download {selectedService?.name}</Title>
            <Vertical gap={12}>
                <Horizontal alignItems="center" gap={12}>
                    <Text>Use Docker?</Text>
                    <Spacer />
                    <SegmentedButtons
                        value={useDocker}
                        onChange={(v) => setUseDocker(v)}
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
                <Horizontal alignItems="center" gap={12}>
                    <Text style={useDocker !== false ? { opacity: 0.4 } : {}}>
                        Download the precompiled release?
                    </Text>
                    <Spacer />
                    <SegmentedButtons
                        value={useReleases}
                        disabled={useDocker !== false}
                        onChange={(v) => setUseReleases(v)}
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
            </Vertical>
            <Horizontal gap={8}>
                <Spacer />
                <Button onClick={closeInstallPopup}>Cancel</Button>
                <Button onClick={install} primary rightSymbol="download">
                    Download
                </Button>
            </Horizontal>
        </Popup>
    );

    const importPopup = (
        <Popup
            show={showImportPopup}
            onDismiss={() => setShowImportPopup(false)}
        >
            <Title>Import from elsewhere</Title>
            <SegmentedButtons
                disabled={importing}
                value={importMethod}
                onChange={(v) => setImportMethod(v)}
                items={[
                    {
                        label: "Git remote",
                        value: "git",
                        rightSymbol: "cloud_download",
                    },
                    {
                        label: "Local storage",
                        value: "localstorage",
                        rightSymbol: "storage",
                    },
                ]}
            />
            {importMethod === "git" && (
                <Input
                    disabled={importing}
                    value={repository}
                    onChange={onRepoChange}
                    label="Repository"
                    placeholder="github.com/user/repo"
                    description="All Git remotes are compatible."
                />
            )}
            {importMethod === "localstorage" && (
                <Input
                    disabled={importing}
                    value={repository}
                    onChange={onRepoChange}
                    label="Service path"
                    description="Absolute path on your local machine where the server is running."
                />
            )}
            {importing && <Progress infinite />}
            <Horizontal gap={12}>
                <Spacer />
                <Button
                    disabled={importing}
                    onClick={() => setShowImportPopup(false)}
                >
                    Cancel
                </Button>
                <Button
                    disabled={importing}
                    onClick={() => installFromElsewhere()}
                    primary
                    rightSymbol={importMethod === "git" ? "download" : "link"}
                >
                    {importMethod === "git" ? "Download" : "Link"}
                </Button>
            </Horizontal>
        </Popup>
    );

    return (
        <Fragment>
            <Header />
            <div className={styles.page}>
                <Horizontal gap={12} alignItems="center">
                    <Title>Marketplace</Title>
                    <Spacer />
                    or
                    <Button
                        onClick={() => setShowImportPopup(true)}
                        rightSymbol="download"
                    >
                        Import from elsewhere
                    </Button>
                </Horizontal>
                {<Error error={error} />}
                <Vertical>
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
            {importPopup}
        </Fragment>
    );
}

import { useFetch } from "../../../hooks/useFetch";
import { BigTitle, Text, Title } from "../../../components/Text/Text";
import { Service as ServiceModel } from "../../../models/service";
import { Fragment, useState } from "react";

import styles from "./InstancesStore.module.sass";
import Service from "../../../components/Service/Service";
import { Horizontal, Vertical } from "../../../components/Layouts/Layouts";
import Spacer from "../../../components/Spacer/Spacer";
import { SegmentedButtons } from "../../../components/SegmentedButton";
import Button from "../../../components/Button/Button";
import Popup from "../../../components/Popup/Popup";
import { InstallMethod, Instances } from "../../../models/instance";
import { api } from "../../../backend/backend";
import { APIError, Errors } from "../../../components/Error/Error";

type Downloading = {
    service: ServiceModel;
};

export default function InstancesStore() {
    const { data: services, error: servicesError } = useFetch<ServiceModel[]>(
        api.services.available.get
    );
    const {
        data: instances,
        error: instancesError,
        reload: reloadInstances,
    } = useFetch<Instances>(api.instances.get);

    const [installMethod, setInstallMethod] = useState<InstallMethod>();

    const [showInstallPopup, setShowInstallPopup] = useState<boolean>(false);

    const [selectedService, setSelectedService] = useState<ServiceModel>();

    const [error, setError] = useState();
    const [popupError, setPopupError] = useState();

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
            // @ts-ignore
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
            .catch(setError)
            .finally(() => {
                setDownloading(
                    downloading.filter(({ service: s }) => s.id !== service.id)
                );
                reloadInstances().catch(console.error);
            });
    };

    const installPopup = (
        <Popup
            show={showInstallPopup}
            onDismiss={() => setShowInstallPopup(false)}
        >
            <Title>Download {selectedService?.name}</Title>
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
            <APIError style={{ margin: 0 }} error={popupError} />
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

    return (
        <Fragment>
            <div className={styles.page}>
                <Horizontal
                    className={styles.title}
                    gap={10}
                    alignItems="center"
                >
                    <BigTitle>Create instance</BigTitle>
                </Horizontal>
                <Errors className={styles.errors}>
                    <APIError error={error} />
                    <APIError error={servicesError} />
                    <APIError error={instancesError} />
                </Errors>
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
                                          ({ service: s }) =>
                                              s.id === service.id
                                      )?.length
                            }
                        />
                    ))}
                </Vertical>
            </div>
            {installPopup}
        </Fragment>
    );
}

import { Fragment, useEffect, useState } from "react";
import { Text, Title } from "../../../../components/Text/Text";
import { Horizontal } from "../../../../components/Layouts/Layouts";
import Spacer from "../../../../components/Spacer/Spacer";
import Button from "../../../../components/Button/Button";
import { useParams } from "react-router-dom";
import useContainer from "../../../../hooks/useContainer";
import Icon from "../../../../components/Icon/Icon";
import ToggleButton from "../../../../components/ToggleButton/ToggleButton";
import Input from "../../../../components/Input/Input";

import styles from "./ContainerSettings.module.sass";
import { api } from "../../../../backend/api/backend";
import { APIError } from "../../../../components/Error/APIError";
import Select, {
    SelectOption,
    SelectValue,
} from "../../../../components/Input/Select";
import VersionTag from "../../../../components/VersionTag/VersionTag";
import classNames from "classnames";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export default function ContainerSettings() {
    const { uuid } = useParams();
    const queryClient = useQueryClient();

    const { container, isLoading: isLoadingContainer } = useContainer(uuid);

    const [displayName, setDisplayName] = useState<string>();
    const [launchOnStartup, setLaunchOnStartup] = useState<boolean>();
    const [version, setVersion] = useState<string>();
    const [versions, setVersions] = useState<string[]>();
    const [versionsLoading, setVersionsLoading] = useState<boolean>(false);

    // undefined = not saved AND never modified
    const [saved, setSaved] = useState<boolean>(undefined);
    const [error, setError] = useState();

    useEffect(() => {
        if (!container) return;
        setLaunchOnStartup(container?.launch_on_startup ?? true);
        setDisplayName(container?.display_name ?? container?.service?.name);
        setVersion(container?.version ?? "latest");
        reloadVersions();
    }, [container]);

    const reloadVersions = (cache = true) => {
        setVersionsLoading(true);
        api.vxContainers
            .container(container.uuid)
            .versions.get(cache)
            .then((data) => {
                setVersions(data?.reverse());
            })
            .catch(setError)
            .finally(() => {
                setVersionsLoading(false);
            });
    };

    const mutationSave = useMutation({
        mutationFn: async () => {
            await api.vxContainers.container(uuid).patch({
                launch_on_startup: launchOnStartup,
                display_name: displayName,
                version: version,
            });
        },
        onSuccess: () => {
            setSaved(true);
        },
        onSettled: () => {
            queryClient.invalidateQueries({
                queryKey: ["containers", uuid],
            });
        },
    });
    const { isLoading: isUploading } = mutationSave;

    const onVersionChange = (v: any) => {
        setVersion(v);
        setSaved(false);
    };

    const versionValue = (
        <SelectValue
            className={classNames({
                [styles.versionValue]: version !== "latest",
            })}
        >
            {version === "latest" ? (
                "Always pull latest version"
            ) : (
                <VersionTag>{version}</VersionTag>
            )}
        </SelectValue>
    );

    return (
        <Fragment>
            <ProgressOverlay
                show={isLoadingContainer || versionsLoading || isUploading}
            />
            <Title className={styles.title}>Settings</Title>
            <APIError error={error} />
            <Input
                label="Container name"
                description="The custom name of your choice for this service"
                value={displayName}
                onChange={(e: any) => {
                    setDisplayName(e.target.value);
                    setSaved(false);
                }}
                disabled={isLoadingContainer}
            />
            <div className={styles.versionSelect}>
                <Select
                    label="Version"
                    onChange={onVersionChange}
                    disabled={isLoadingContainer || versionsLoading}
                    // @ts-ignore
                    value={versionValue}
                >
                    {versions?.includes("latest") && (
                        <SelectOption value="latest">
                            Always pull latest version
                        </SelectOption>
                    )}
                    {versions?.map((v) => {
                        if (v === "latest") {
                            return null;
                        }
                        return (
                            <SelectOption
                                key={v}
                                value={v}
                                className={styles.versionOption}
                            >
                                <VersionTag>{v}</VersionTag>
                            </SelectOption>
                        );
                    })}
                </Select>
                <Button
                    rightIcon="refresh"
                    onClick={() => reloadVersions(false)}
                    disabled={isLoadingContainer || versionsLoading}
                >
                    Refresh
                </Button>
            </div>
            <Horizontal className={styles.toggle} alignItems="center">
                <Text>Launch on Startup</Text>
                <Spacer />
                <ToggleButton
                    value={launchOnStartup}
                    onChange={(v) => {
                        setLaunchOnStartup(v);
                        setSaved(false);
                    }}
                    disabled={isLoadingContainer}
                />
            </Horizontal>
            <Button
                primary
                large
                onClick={async () => mutationSave.mutate()}
                rightIcon="save"
                loading={isUploading}
                disabled={saved || saved === undefined}
            >
                Save
            </Button>
            {!isUploading && saved && (
                <Horizontal
                    className={styles.saved}
                    alignItems="center"
                    gap={4}
                >
                    <Icon name="check" />
                    Saved!
                </Horizontal>
            )}
        </Fragment>
    );
}

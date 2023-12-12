import { useEffect, useState } from "react";
import { Horizontal } from "../../../../components/Layouts/Layouts";
import Spacer from "../../../../components/Spacer/Spacer";
import {
    Button,
    MaterialIcon,
    Paragraph,
    SelectField,
    SelectOption,
    TextField,
    Title,
} from "@vertex-center/components";
import { useParams } from "react-router-dom";
import useContainer from "../../hooks/useContainer";
import ToggleButton from "../../../../components/ToggleButton/ToggleButton";
import styles from "./ContainerSettings.module.sass";
import { APIError } from "../../../../components/Error/APIError";
import VersionTag from "../../../../components/VersionTag/VersionTag";
import classNames from "classnames";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import Content from "../../../../components/Content/Content";
import { API } from "../../backend/api";

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
        setLaunchOnStartup(container.launch_on_startup);
        setDisplayName(container.name);
        setVersion(container?.image_tag ?? "latest");
        reloadVersions();
    }, [container]);

    const reloadVersions = (cache = true) => {
        setVersionsLoading(true);
        API.getVersions(container.id, cache)
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
            await API.patchContainer(uuid, {
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
        <div
            className={classNames({
                [styles.versionValue]: version !== "latest",
            })}
        >
            {version === "latest" ? (
                "Always pull latest version"
            ) : (
                <VersionTag>{version}</VersionTag>
            )}
        </div>
    );

    return (
        <Content>
            <Title variant="h2">Settings</Title>
            <ProgressOverlay
                show={isLoadingContainer || versionsLoading || isUploading}
            />
            <APIError error={error} />
            <Horizontal alignItems="center">
                <Paragraph>Launch on Startup</Paragraph>
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
            <TextField
                id="container-name"
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
                <SelectField
                    id="container-version"
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
                            <SelectOption key={v} value={v}>
                                <VersionTag>{v}</VersionTag>
                            </SelectOption>
                        );
                    })}
                </SelectField>
                <Button
                    rightIcon={<MaterialIcon icon="refresh" />}
                    onClick={() => reloadVersions(false)}
                    disabled={isLoadingContainer || versionsLoading}
                >
                    Refresh
                </Button>
            </div>
            <Horizontal justifyContent="flex-end">
                {!isUploading && saved && (
                    <Horizontal
                        className={styles.saved}
                        alignItems="center"
                        gap={4}
                    >
                        <MaterialIcon icon="check" />
                        Saved!
                    </Horizontal>
                )}
                <Button
                    variant="colored"
                    onClick={async () => mutationSave.mutate()}
                    rightIcon={<MaterialIcon icon="save" />}
                    disabled={isUploading || saved || saved === undefined}
                >
                    Save
                </Button>
            </Horizontal>
        </Content>
    );
}

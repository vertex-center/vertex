import { Fragment, useEffect, useState } from "react";
import { Button, Input, MaterialIcon, Title } from "@vertex-center/components";
import { Horizontal } from "../../../components/Layouts/Layouts";
import { APIError } from "../../../components/Error/APIError";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import Content from "../../../components/Content/Content";
import { useSettings } from "../hooks/useSettings";
import { usePatchSettings } from "../hooks/usePatchSettings";

export default function SettingsNotifications() {
    const [webhook, setWebhook] = useState<string>();
    const [changed, setChanged] = useState(false);

    const { settings, errorSettings, isLoadingSettings } = useSettings();

    useEffect(() => {
        setWebhook(settings?.webhook);
    }, [settings]);

    const onWebhookChange = (e: any) => {
        setWebhook(e.target.value);
        setChanged(true);
    };

    const { patchSettings, isPatchingSettings, errorPatchingSettings } =
        usePatchSettings({
            onSuccess: () => setChanged(false),
        });

    const error = errorSettings || errorPatchingSettings;
    const isLoading = isLoadingSettings;

    return (
        <Content>
            <Title variant="h2">Notifications</Title>
            <ProgressOverlay show={isLoading || isPatchingSettings} />
            <APIError error={error} />
            {!error && (
                <Fragment>
                    <Input
                        id="webhook"
                        label="Webhook"
                        value={webhook}
                        onChange={onWebhookChange}
                        disabled={isLoading}
                        placeholder={isLoading && "Loading..."}
                    />
                    <Horizontal
                        gap={20}
                        justifyContent="flex-end"
                        style={{ marginTop: 15 }}
                    >
                        <Button
                            variant="colored"
                            rightIcon={<MaterialIcon icon="save" />}
                            onClick={() => patchSettings({ webhook })}
                            disabled={!changed || isPatchingSettings}
                        >
                            Save
                        </Button>
                    </Horizontal>
                </Fragment>
            )}
        </Content>
    );
}

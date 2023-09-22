import { Fragment, useEffect, useState } from "react";
import Input from "../../components/Input/Input";
import { Title } from "../../components/Text/Text";
import Button from "../../components/Button/Button";
import Loading from "../../components/Loading/Loading";
import { useFetch } from "../../hooks/useFetch";
import { api } from "../../backend/backend";

import styles from "./SettingsNotifications.module.sass";
import { Vertical } from "../../components/Layouts/Layouts";
import { APIError } from "../../components/Error/Error";

export default function SettingsNotifications() {
    const [webhook, setWebhook] = useState<string>();
    const [changed, setChanged] = useState(false);
    const [saving, setSaving] = useState(false);

    const { data: settings, error } = useFetch<Settings>(api.settings.get);

    useEffect(() => {
        setWebhook(settings?.notifications?.webhook);
    }, [settings]);

    const onWebhookChange = (e: any) => {
        setWebhook(e.target.value);
        setChanged(true);
    };

    const onSave = () => {
        setSaving(true);
        api.settings
            .patch({ notifications: { webhook } })
            .then(() => setChanged(false))
            .catch(console.error)
            .finally(() => setSaving(false));
    };

    return (
        <Vertical gap={20}>
            <Title className={styles.title}>Notifications</Title>
            <APIError error={error} />
            {!error && (
                <Fragment>
                    <Input
                        label="Webhook"
                        value={webhook}
                        onChange={onWebhookChange}
                    />
                    <Button
                        large
                        rightSymbol="save"
                        onClick={onSave}
                        disabled={!changed || saving}
                    >
                        Save
                    </Button>
                    {saving && <Loading />}
                </Fragment>
            )}
        </Vertical>
    );
}

import {
    Button,
    Horizontal,
    List,
    ListDescription,
    ListIcon,
    ListInfo,
    ListItem,
    ListTitle,
    Title,
} from "@vertex-center/components";
import Content from "../../../components/Content/Content";
import { useServerEvent } from "../../../hooks/useEvent";
import { useState } from "react";
import NoItems from "../../../components/NoItems/NoItems";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { Check, ListChecks, Play, WarningCircle } from "@phosphor-icons/react";

export default function SettingsChecks() {
    const [isChecking, setIsChecking] = useState(false);

    const [checks, setChecks] = useState({});

    useServerEvent(
        // @ts-ignore
        window.api_urls.admin,
        "/admin/checks",
        {
            check: (e) => {
                const d = JSON.parse(e.data);
                setChecks((c) => ({ ...c, [d.id]: d }));
            },
            done: () => setIsChecking(false),
        },
        !isChecking
    );

    const runChecks = () => {
        setIsChecking(true);
        setChecks({});
    };

    return (
        <Content>
            <ProgressOverlay show={isChecking} />
            <Title variant="h2">Checks</Title>
            <Horizontal>
                <Button
                    variant="colored"
                    rightIcon={<Play />}
                    onClick={runChecks}
                    disabled={isChecking}
                >
                    Run checks
                </Button>
            </Horizontal>
            {checks && Object.keys(checks).length === 0 && !isChecking && (
                <NoItems
                    icon={<ListChecks />}
                    text="Run checks to see if there are any issues with your installation."
                />
            )}
            <List>
                {Object.values(checks ?? {}).map((check: any) => (
                    <ListItem key={check?.id}>
                        <ListIcon>
                            {check?.error === "" ? (
                                <Check />
                            ) : (
                                <WarningCircle />
                            )}
                        </ListIcon>
                        <ListInfo>
                            <ListTitle>{check?.name}</ListTitle>
                            <ListDescription>
                                {check?.error === "" ? "OK" : check?.error}
                            </ListDescription>
                        </ListInfo>
                    </ListItem>
                ))}
            </List>
        </Content>
    );
}

import {
    Button,
    Horizontal,
    List,
    ListDescription,
    ListIcon,
    ListInfo,
    ListItem,
    ListTitle,
    MaterialIcon,
    Title,
} from "@vertex-center/components";
import Content from "../../../components/Content/Content";
import { useServerEvent } from "../../../hooks/useEvent";
import { useState } from "react";
import NoItems from "../../../components/NoItems/NoItems";
import { ProgressOverlay } from "../../../components/Progress/Progress";

export default function SettingsChecks() {
    const [isChecking, setIsChecking] = useState(false);

    const [checks, setChecks] = useState({});

    useServerEvent(
        "7500",
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
                    rightIcon={<MaterialIcon icon="play_arrow" />}
                    onClick={runChecks}
                    disabled={isChecking}
                >
                    Run checks
                </Button>
            </Horizontal>
            {checks && Object.keys(checks).length === 0 && !isChecking && (
                <NoItems
                    icon="checklist"
                    text="Run checks to see if there are any issues with your installation."
                />
            )}
            <List>
                {Object.values(checks ?? {}).map((check: any) => (
                    <ListItem key={check?.id}>
                        <ListIcon>
                            <MaterialIcon
                                icon={check?.error !== "" ? "error" : "check"}
                            />
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

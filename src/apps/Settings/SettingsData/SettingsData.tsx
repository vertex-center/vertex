import Content from "../../../components/Content/Content";
import {
    Button,
    Horizontal,
    List,
    ListActions,
    ListDescription,
    ListIcon,
    ListInfo,
    ListItem,
    ListTitle,
    MaterialIcon,
    Paragraph,
    Title,
    Vertical,
} from "@vertex-center/components";
import { SiPostgresql, SiSqlite } from "@icons-pack/react-simple-icons";
import styles from "./SettingsData.module.sass";
import Popup from "../../../components/Popup/Popup";
import { Fragment, useState } from "react";
import { useDBMS } from "../hooks/useDBMS";
import Progress, {
    ProgressOverlay,
} from "../../../components/Progress/Progress";
import { APIError } from "../../../components/Error/APIError";
import { useDBMSMutation } from "../hooks/useDBMSMutation";
import { useQueryClient } from "@tanstack/react-query";

type DatabaseProps = {
    name: string;
    icon: JSX.Element;
    title: string;
    description: string;
    installed: boolean;
    onMigrate?: (db: string) => void;
    hideActions?: boolean;
};

function Database(props: Readonly<DatabaseProps>) {
    const {
        name,
        icon,
        title,
        description,
        installed,
        onMigrate,
        hideActions,
    } = props;

    let actions = null;
    if (installed) {
        actions = (
            <Horizontal className={styles.tag} alignItems="center">
                <MaterialIcon icon="check" />
                Active
            </Horizontal>
        );
    } else if (!hideActions) {
        actions = (
            <Button
                variant="danger"
                rightIcon={<MaterialIcon icon="restart_alt" />}
                onClick={() => onMigrate(name)}
            >
                Migrate
            </Button>
        );
    }

    return (
        <ListItem>
            <ListIcon>{icon}</ListIcon>
            <ListInfo>
                <ListTitle>{title}</ListTitle>
                <ListDescription>{description}</ListDescription>
            </ListInfo>
            <ListActions>{actions}</ListActions>
        </ListItem>
    );
}

export default function SettingsData() {
    const queryClient = useQueryClient();

    const [showPopup, setShowPopup] = useState(false);
    const [selectedDB, setSelectedDB] = useState<string>();

    const { dbms, isLoadingDbms, errorDbms } = useDBMS();
    const { migrate, isMigrating, errorMigrate } = useDBMSMutation({
        onSettled: () => {
            setShowPopup(false);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({
                queryKey: ["admin_data_dbms"],
            });
        },
    });

    const askMigrate = (db: string) => {
        setSelectedDB(db);
        setShowPopup(true);
    };

    const dismissPopup = () => setShowPopup(false);

    const askMigrateActions = (
        <Fragment>
            <Button variant="outlined" onClick={dismissPopup}>
                Cancel
            </Button>
            <Button
                variant="danger"
                rightIcon={<MaterialIcon icon="restart_alt" />}
                onClick={() => migrate(selectedDB)}
            >
                Migrate
            </Button>
        </Fragment>
    );

    return (
        <Content>
            <Title variant="h2">Data</Title>
            <Paragraph>
                You can choose between <b>SQLite</b> and <b>Postgres</b> to
                store your Vertex data. You don't need to worry about installing
                or configuring the database, Vertex will do that for you.
                {/*Vertex data can include users, permissions, settings....*/}
            </Paragraph>
            <ProgressOverlay show={isLoadingDbms} />
            <APIError error={errorDbms || errorMigrate} />
            <List>
                <Database
                    name="sqlite"
                    icon={<SiSqlite />}
                    title="SQLite"
                    description="Recommended for small setups."
                    installed={dbms === "sqlite"}
                    onMigrate={askMigrate}
                    hideActions={isMigrating || isLoadingDbms}
                />
                <Database
                    name="postgres"
                    icon={<SiPostgresql />}
                    title="Postgres"
                    description="Recommended for larger installations."
                    installed={dbms === "postgres"}
                    onMigrate={askMigrate}
                    hideActions={isMigrating || isLoadingDbms}
                />
            </List>
            {isMigrating && (
                <Vertical gap={12}>
                    <Paragraph>Migration to {selectedDB} ongoing...</Paragraph>
                    <Progress infinite />
                </Vertical>
            )}
            <Popup
                show={showPopup}
                onDismiss={dismissPopup}
                title="Migrate database?"
                actions={askMigrateActions}
            >
                <Paragraph>
                    Are you sure you want to migrate from {dbms} to {selectedDB}
                    ?
                </Paragraph>
            </Popup>
        </Content>
    );
}

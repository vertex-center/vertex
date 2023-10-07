import { PropsWithChildren, useState } from "react";

import styles from "./Update.module.sass";
import Button from "../Button/Button";
import Icon from "../Icon/Icon";
import Progress from "../Progress";
import ListItem from "../List/ListItem";
import List from "../List/List";
import ListInfo from "../List/ListInfo";
import ListTitle from "../List/ListTitle";
import ListDescription from "../List/ListDescription";
import ListActions from "../List/ListActions";

export function Updates(props: Readonly<PropsWithChildren>) {
    if (!props.children) return null;
    return <List {...props} />;
}

type Props = {
    name: string;
    version?: string;
    available?: boolean;

    current_version?: string;
    latest_version?: string;

    onUpdate: () => Promise<void>;
};

export default function Update(props: Readonly<Props>) {
    const { name, available, version, latest_version } = props;

    const [isLoading, setIsLoading] = useState(false);

    const onUpdate = () => {
        setIsLoading(true);
        props.onUpdate().finally(() => setIsLoading(false));
    };

    let status: JSX.Element;
    if (available) {
        if (isLoading) {
            status = <Progress infinite />;
        } else {
            status = (
                <ListActions>
                    <Button rightIcon="download" onClick={onUpdate}>
                        Update
                    </Button>
                </ListActions>
            );
        }
    } else {
        status = (
            <div className={styles.status}>
                <Icon name="check" />
                Up-to-date
            </div>
        );
    }

    return (
        <ListItem>
            <ListInfo>
                <ListTitle className={styles.name}>{name}</ListTitle>
                {version && (
                    <ListDescription>
                        <code className={styles.version}>
                            {version}
                            {latest_version && " -> " + latest_version}
                        </code>
                    </ListDescription>
                )}
            </ListInfo>
            {status}
        </ListItem>
    );
}

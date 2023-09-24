import { PropsWithChildren, useState } from "react";

import styles from "./Update.module.sass";
import Button from "../Button/Button";
import Symbol from "../Symbol/Symbol";
import Progress from "../Progress";
import ListItem from "../List/ListItem";
import List from "../List/List";
import ListInfo from "../List/ListInfo";
import ListTitle from "../List/ListTitle";
import ListDescription from "../List/ListDescription";
import ListActions from "../List/ListActions";

export function Updates(props: PropsWithChildren) {
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

export default function Update(props: Props) {
    const { name, available, version, latest_version } = props;

    const [isLoading, setIsLoading] = useState(false);

    const onUpdate = () => {
        setIsLoading(true);
        props.onUpdate().finally(() => setIsLoading(false));
    };

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
            {available ? (
                isLoading ? (
                    <Progress infinite />
                ) : (
                    <ListActions>
                        <Button rightSymbol="download" onClick={onUpdate}>
                            Update
                        </Button>
                    </ListActions>
                )
            ) : (
                <div className={styles.status}>
                    <Symbol name="check" />
                    Up-to-date
                </div>
            )}
        </ListItem>
    );
}

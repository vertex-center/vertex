import Content from "../../../../components/Content/Content";
import {
    List,
    ListDescription,
    ListIcon,
    ListInfo,
    ListItem,
    ListTitle,
    MaterialIcon,
    Paragraph,
    Title,
} from "@vertex-center/components";
import { useCredentials } from "../../hooks/useCredentials";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";

export default function AccountSecurity() {
    const { credentials, isLoadingCredentials, errorCredentials } =
        useCredentials();

    const isLoading = isLoadingCredentials;
    const error = errorCredentials;

    return (
        <Content>
            <Title variant="h2">Security</Title>
            <ProgressOverlay show={isLoading} />
            <APIError error={error} />
            <Paragraph>
                This page shows all the methods currently enabled to
                authenticate with your account.
            </Paragraph>
            <List>
                {credentials?.map((cred, i) => (
                    <ListItem key={i}>
                        <ListIcon>
                            <MaterialIcon icon="password" />
                        </ListIcon>
                        <ListInfo>
                            <ListTitle>{cred.name}</ListTitle>
                            <ListDescription>
                                {cred.description}
                            </ListDescription>
                        </ListInfo>
                    </ListItem>
                ))}
            </List>
        </Content>
    );
}

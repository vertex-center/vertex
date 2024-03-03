import Content from "../../../../components/Content/Content";
import {
    Button,
    FormItem,
    Horizontal,
    Input,
    Title,
    Vertical,
} from "@vertex-center/components";
import useUser from "../../hooks/useUser";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import Spacer from "../../../../components/Spacer/Spacer";
import { usePatchUser } from "../../hooks/usePatchUser";
import { ChangeEvent, useEffect, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import Saved from "../../../../components/Saved/Saved";
import { FloppyDiskBack } from "@phosphor-icons/react";

export default function AccountInfo() {
    const queryClient = useQueryClient();
    const { user, isLoadingUser, errorUser } = useUser();

    const [saved, setSaved] = useState<boolean>(undefined);

    const [username, setUsername] = useState(user?.username);

    useEffect(() => {
        setUsername(user?.username);
    }, [user]);

    const { patchUser, isPatchingUser, errorPatchUser } = usePatchUser({
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["user"] });
            setSaved(true);
        },
    });

    const isLoading = isPatchingUser || isLoadingUser;
    const error = errorPatchUser || errorUser;

    const save = () => {
        patchUser({
            username,
        });
    };

    const onChangeUsername = (e: ChangeEvent<HTMLInputElement>) => {
        setUsername(e.target.value);
        setSaved(false);
    };

    return (
        <Content>
            <Title variant="h2">Information</Title>
            <ProgressOverlay show={isLoading} />
            <APIError error={error} />
            <Vertical gap={20}>
                <FormItem label="Username">
                    <Input
                        value={username}
                        onChange={onChangeUsername}
                        disabled={isLoading}
                        required
                    />
                </FormItem>
                <Horizontal gap={20}>
                    <Spacer />
                    <Saved show={saved} />
                    <Button
                        variant="colored"
                        rightIcon={<FloppyDiskBack />}
                        onClick={save}
                        disabled={
                            isLoading || saved === true || saved === undefined
                        }
                    >
                        Save
                    </Button>
                </Horizontal>
            </Vertical>
        </Content>
    );
}

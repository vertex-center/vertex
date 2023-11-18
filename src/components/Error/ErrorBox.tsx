import { Box, BoxProps } from "@vertex-center/components";

type Props = BoxProps & {
    error?: any;
};

export default function ErrorBox(props: Readonly<Props>) {
    const { error, className, ...others } = props;
    let err = error?.message ?? "An unknown error has occurred.";
    return (
        <Box type="error" {...others}>
            {err}
        </Box>
    );
}

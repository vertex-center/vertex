import { Box, BoxProps, Paragraph } from "@vertex-center/components";

type Props = Omit<BoxProps, "type"> & {
    error?: any;
};

export default function ErrorBox(props: Readonly<Props>) {
    const { error, className, ...others } = props;
    let err = error?.message ?? "An unknown error has occurred.";
    return (
        <Box type="error" {...others}>
            <Paragraph>{err}</Paragraph>
        </Box>
    );
}

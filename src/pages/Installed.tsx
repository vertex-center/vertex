function Application() {
    return (
        <div className="flex flex-col rounded-md bg-zinc-100 px-4 py-2">
            <h3 className="text-lg font-medium">Vertex Spotify</h3>
        </div>
    );
}

export default function Installed() {
    return (
        <div className="flex flex-col border-separate border-amber-500 gap-4 m-4">
            <Application />
            <Application />
            <Application />
        </div>
    );
}

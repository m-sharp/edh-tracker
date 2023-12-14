import { usePathname } from 'next/navigation'

export const dynamicParams = false;

// Instead of prefetching from an API, for now use the next dynamic routing basically to just pass around IDs
// ToDo: Eventually, split out the API web server from the frontend entirely and make use of all the fancy Next.js stuff
export function generateStaticParams() {
    // Generate no routes for the build
    // ToDo: This just doesn't work with an empty array...Need to figure out a different way to pass contextual data between react pages besides route params
    // return [{ id: '1' }, { id: '2' }, { id: '3' }];
    return [{}];
}

export default function Page({ params }) {
    const pathname = usePathname();

    return (
        <main>
            <h1>Hello, Player {params.id}'s Page!</h1>
            <p>Path name: {pathname}</p>
        </main>
    )
}

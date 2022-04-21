export function AboutPage() {
  return (
    <main className="flex flex-col items-center">
      <p className="text-2xl my-16">About</p>
      <p className="text-xl">
        A secure link against man-in-the-middle traffic analysis attacks
      </p>
      <p className="text-xl">
        based on quic, one uuid represent a unique connection
      </p>
      <a
        className="text-xl my-16"
        href="https://github.com/wkj9893/masky"
        target="_blank"
      >
        source code
      </a>
    </main>
  );
}

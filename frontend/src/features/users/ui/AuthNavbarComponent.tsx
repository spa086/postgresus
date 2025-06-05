export function AuthNavbarComponent() {
  return (
    <div className="m-3 flex h-[65px] items-center p-3">
      <div className="flex items-center gap-3 hover:opacity-80">
        <a href="https://postgresus.com" target="_blank" rel="noreferrer">
          <img className="h-[35px] w-[35px]" src="/logo.svg" />
        </a>

        <div className="text-xl font-bold">
          <a href="https://postgresus.com" className="text-black" target="_blank" rel="noreferrer">
            Postgresus
          </a>
        </div>
      </div>

      <div className="mr-3 ml-auto flex gap-5">
        <a
          className="hover:opacity-80"
          href="https://postgresus.com/community"
          target="_blank"
          rel="noreferrer"
        >
          Community
        </a>

        <a
          className="hover:opacity-80"
          href="https://github.com/postgresus/postgresus"
          target="_blank"
          rel="noreferrer"
        >
          GitHub
        </a>
      </div>
    </div>
  );
}

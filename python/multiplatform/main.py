import sys
import anyio
import dagger

async def test():
  platforms = ["linux/amd64", "linux/arm64"]

  async with dagger.Connection(dagger.Config(log_output=sys.stderr)) as client:

    # get reference to the local project
    src_id = await client.host().directory(".").id()

    variants = []
    for platform in platforms:

      python = (
        client.container(platform=platform)
        .from_(f"python:3.10-slim-buster")
      )

      multistage = await python.rootfs().with_directory("/src", src_id).id()

      python = (
        python.with_rootfs(multistage)
        .with_workdir("/src")
        .with_exec(["pip", "install", "dagger-io"])
        .with_entrypoint(["python", "/src/main.py"])
      )

      container_id = await python.id()
      variants.append(container_id)

    await client.container().publish("kylepenfound/hello-python:latest", variants)
    print("All tasks have finished")

if __name__ == "__main__":
    anyio.run(test)
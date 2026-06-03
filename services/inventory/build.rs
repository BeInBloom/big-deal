fn main() -> Result<(), Box<dyn std::error::Error>> {
    println!("cargo:rerun-if-changed=../../shared/proto/inventory/v1/inventory.proto");

    tonic_build::configure()
        .build_server(true)
        .build_client(false)
        .compile_protos(
            &["../../shared/proto/inventory/v1/inventory.proto"],
            &["../../shared/proto"],
        )?;

    Ok(())
}

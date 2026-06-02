fn main() -> Result<(), Box<dyn std::error::Error>> {
    tonic_build::configure()
        .build_server(true)
        .build_client(false)
        .compile_protos(
            &["../../shared/proto/payment/v1/payment.proto"],
            &["../../shared/proto"],
        )?;

    Ok(())
}

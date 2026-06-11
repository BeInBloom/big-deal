use crate::bson;

pub trait IntoBson {
    fn into_bson(self) -> bson::Bson;
}

impl IntoBson for String {
    fn into_bson(self) -> bson::Bson {
        bson::Bson::String(self)
    }
}

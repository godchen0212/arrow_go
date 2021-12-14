#include <arrow/api.h>
#include<arrow/memory_pool.h>
#include <arrow/io/api.h>
#include <parquet/arrow/reader.h>
#include <parquet/arrow/writer.h>
#include <parquet/exception.h>
#include <unistd.h>
#include <iostream>

std::shared_ptr<arrow::Table> generate_table() {
    arrow::Int64Builder i64builder;
    PARQUET_THROW_NOT_OK(i64builder.AppendValues({1, 2, 3, 4, 5}));
    auto Size = 320000 * 5;
    for (int i = 0; i < Size; i++) {
        i64builder.Append(i);
    }
    std::shared_ptr<arrow::Array> i64array;
    PARQUET_THROW_NOT_OK(i64builder.Finish(&i64array));

    std::shared_ptr<arrow::Schema> schema = arrow::schema(
            {arrow::field("int", arrow::int64())});

    return arrow::Table::Make(schema, {i64array});
}


void write_parquet_file(const arrow::Table &table) {
    auto mempool = arrow::default_memory_pool();
    std::shared_ptr<arrow::io::FileOutputStream> outfile;

    std::cout << "write::bytes-1::" << mempool->bytes_allocated() << std::endl;
    PARQUET_ASSIGN_OR_THROW(
            outfile, arrow::io::FileOutputStream::Open("parquet-arrow-example2.parquet"));
    PARQUET_THROW_NOT_OK(
            //parquet::arrow::WriteTable(table, arrow::default_memory_pool(), outfile, 3));
            parquet::arrow::WriteTable(table, mempool, outfile, 64*1024*1024));

    std::cout << "write::bytes-2::" << mempool->bytes_allocated() << std::endl;

}

void read_whole_file() {

    auto mempool2 = arrow::default_memory_pool();
    std::cout << "address2:" << mempool2 << std::endl;

    auto mempool = arrow::MemoryPool::CreateDefault();
    std::cout << "address:" << mempool.get() << std::endl;
    // auto mempool2 = mempool->CreateDefault();
    std::cout << "Reading parquet-arrow-example2.parquet at once" << std::endl;
    std::shared_ptr<arrow::io::ReadableFile> infile;
    PARQUET_ASSIGN_OR_THROW(infile,
                            arrow::io::ReadableFile::Open("parquet-arrow-example2.parquet",
                                                          mempool.get()));

    std::cout << "read::bytes-1::" << mempool->bytes_allocated() << std::endl;
    std::unique_ptr<parquet::arrow::FileReader> reader;
    PARQUET_THROW_NOT_OK(
            //parquet::arrow::OpenFile(infile, arrow::default_memory_pool(), &reader));
            parquet::arrow::OpenFile(infile, mempool.get(), &reader));
    std::shared_ptr<arrow::Table> table;
    PARQUET_THROW_NOT_OK(reader->ReadTable(&table));
    std::cout << "Loaded " << table->num_rows() << " rows in " << table->num_columns()
              << " columns." << std::endl;

    std::cout << "read::bytes-2::" << mempool->bytes_allocated() << std::endl;
    mempool->ReleaseUnused();
//  mempool = mempool->CreateDefault();
    std::cout << "read::bytes-after callReleaseUnused::" << mempool->bytes_allocated() << std::endl;

}

int main() {
    //std::shared_ptr<arrow::Table> table = generate_table();
    //write_parquet_file(*table);
    //std::cout<<"write parquet_file done"<<std::endl;
    for (int i = 0; i < 100; i++) {
        std::cout << "start " << std::endl;
        read_whole_file();
        std::cout << "end " << std::endl;
        sleep(2);
    }
    std::cout << "start sleep 100" << std::endl;
    sleep(100);
}

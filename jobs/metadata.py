import pyspark.sql.functions as F
import pyspark.sql.types as T
from pyspark.sql import SparkSession
from pyspark import SparkConf

FOS = {"civil engineering": (410, 45, 4), "photochemistry": (304, 31, 3), "nuclear chemistry": (303, 31, 3),
       "information retrieval": (409, 44, 4), "inorganic chemistry": (301, 31, 3), "pure mathematics": (312, 33, 3),
       "demography": (111, 12, 1), "climatology": (313, 32, 3), "soil science": (207, 23, 2),
       "political science": (111, 12, 1), "engineering drawing": (410, 45, 4), "theology": (107, 11, 1),
       "public economics": (112, 12, 1), "oceanography": (313, 34, 3), "aeronautics": (402, 41, 4),
       "agricultural engineering": (207, 23, 2), "advertising": (112, 12, 1), "religious studies": (106, 11, 1),
       "chemical engineering": (403, 42, 4), "development economics": (112, 12, 1), "ancient history": (101, 11, 1),
       "geography": (317, 34, 3), "environmental protection": (410, 45, 4), "economic geography": (317, 34, 3),
       "marketing": (112, 12, 1), "topology": (312, 33, 3), "epistemology": (108, 11, 1),
       "computer architecture": (409, 44, 4), "biophysics": (201, 21, 2), "archaeology": (101, 11, 1),
       "art history": (103, 11, 1), "mathematical analysis": (312, 33, 3), "psychoanalysis": (110, 12, 1),
       "materials science": (406, 43, 4), "psychology": (110, 12, 1), "computer vision": (409, 44, 4),
       "nursing": (205, 22, 2), "agronomy": (207, 23, 2), "mining engineering": (410, 45, 4),
       "pulp and paper industry": (405, 43, 4), "particle physics": (309, 32, 3), "law": (113, 12, 1),
       "molecular physics": (308, 32, 3), "financial economics": (112, 12, 1), "mathematics education": (312, 33, 3),
       "engineering ethics": (108, 11, 1), "computational biology": (201, 21, 2), "geophysics": (315, 34, 3),
       "aesthetics": (108, 11, 1), "distributed computing": (409, 44, 4), "psychotherapist": (110, 12, 1),
       "control engineering": (407, 43, 4), "financial system": (112, 12, 1), "philosophy": (108, 11, 1),
       "industrial engineering": (401, 41, 4), "biological engineering": (403, 42, 4), "microeconomics": (112, 12, 1),
       "genealogy": (102, 11, 1), "speech recognition": (409, 44, 4), "natural resource economics": (112, 12, 1),
       "computational science": (409, 44, 4), "economics": (112, 12, 1), "history": (102, 11, 1),
       "gender studies": (111, 12, 1), "business": (112, 12, 1), "electronic engineering": (408, 44, 4),
       "bioinformatics": (201, 21, 2), "biochemical engineering": (403, 42, 4), "microbiology": (204, 22, 2),
       "fishery": (207, 23, 2), "quantum electrodynamics": (309, 32, 3), "discrete mathematics": (312, 33, 3),
       "geochemistry": (316, 34, 3), "cancer research": (205, 22, 2), "combinatorial chemistry": (304, 31, 3),
       "theoretical physics": (307, 32, 3), "anthropology": (106, 11, 1), "hydrology": (318, 34, 3),
       "regional science": (111, 12, 1), "remote sensing": (315, 34, 3), "neoclassical economics": (112, 12, 1),
       "performance art": (103, 11, 1), "classical mechanics": (402, 41, 4), "ethnology": (106, 11, 1),
       "computer engineering": (408, 44, 4), "stereochemistry": (301, 31, 3), "waste management": (410, 45, 4),
       "manufacturing engineering": (401, 41, 4), "chemistry": (301, 31, 3), "earth science": (314, 34, 3),
       "mechanics": (402, 41, 4), "economic growth": (112, 12, 1), "pedagogy": (109, 12, 1),
       "environmental ethics": (108, 11, 1), "algorithm": (409, 44, 4), "agroforestry": (207, 23, 2),
       "industrial organization": (401, 41, 4), "data science": (409, 44, 4), "simulation": (409, 44, 4),
       "biological system": (201, 21, 2), "control theory": (401, 41, 4), "animal science": (203, 21, 2),
       "petroleum engineering": (403, 42, 4), "classical economics": (112, 12, 1), "biochemistry": (201, 21, 2),
       "engineering management": (401, 41, 4), "traditional medicine": (205, 22, 2),
       "architectural engineering": (410, 45, 4), "sociology": (111, 12, 1), "clinical psychology": (110, 12, 1),
       "astronomy": (311, 32, 3), "combinatorics": (312, 33, 3), "arithmetic": (312, 33, 3), "commerce": (112, 12, 1),
       "biology": (201, 21, 2), "environmental planning": (410, 45, 4), "horticulture": (207, 23, 2),
       "chemical physics": (303, 31, 3), "radiochemistry": (303, 31, 3), "condensed matter physics": (307, 32, 3),
       "international economics": (112, 12, 1), "transport engineering": (410, 45, 4), "socioeconomics": (112, 12, 1),
       "statistics": (312, 33, 3), "classics": (103, 11, 1), "geology": (314, 34, 3),
       "mathematical economics": (112, 12, 1), "paleontology": (314, 34, 3), "astrobiology": (201, 21, 2),
       "operations management": (112, 12, 1), "water resource management": (318, 34, 3),
       "geotechnical engineering": (410, 45, 4), "keynesian economics": (112, 12, 1), "cognitive science": (206, 22, 2),
       "forestry": (207, 23, 2), "zoology": (203, 21, 2), "public administration": (112, 12, 1),
       "polymer chemistry": (306, 31, 3), "composite material": (405, 43, 4), "parallel computing": (409, 44, 4),
       "nuclear magnetic resonance": (309, 32, 3), "environmental health": (205, 22, 2), "management": (112, 12, 1),
       "applied psychology": (110, 12, 1), "geodesy": (315, 34, 3), "cell biology": (201, 21, 2),
       "systems engineering": (407, 44, 4), "neuroscience": (206, 22, 2), "anesthesia": (205, 22, 2),
       "computational physics": (310, 32, 3), "mathematics": (312, 33, 3), "optometry": (205, 22, 2),
       "acoustics": (402, 41, 4), "algebra": (312, 33, 3), "food science": (305, 31, 3),
       "agricultural economics": (207, 23, 2), "visual arts": (103, 11, 1), "atomic physics": (309, 32, 3),
       "physiology": (205, 22, 2), "cartography": (315, 34, 3), "mineralogy": (316, 34, 3),
       "engineering physics": (402, 41, 4), "law and economics": (112, 12, 1), "econometrics": (112, 12, 1),
       "mathematical optimization": (312, 33, 3), "gerontology": (205, 22, 2), "geomorphology": (315, 34, 3),
       "anatomy": (204, 22, 2), "actuarial science": (112, 12, 1), "risk analysis": (112, 12, 1),
       "operations research": (112, 12, 1), "welfare economics": (112, 12, 1), "automotive engineering": (402, 31, 3),
       "organic chemistry": (301, 31, 3), "environmental economics": (112, 12, 1), "labour economics": (112, 12, 1),
       "pharmacology": (205, 22, 2), "atmospheric sciences": (313, 34, 3), "environmental engineering": (410, 45, 4),
       "media studies": (103, 11, 1), "applied mathematics": (312, 33, 3), "social psychology": (110, 12, 1),
       "physical chemistry": (303, 31, 3), "market economy": (112, 12, 1), "developmental psychology": (110, 12, 1),
       "nanotechnology": (406, 43, 4), "molecular biology": (201, 21, 2), "polymer science": (306, 31, 3),
       "environmental chemistry": (305, 31, 3), "metallurgy": (405, 43, 4), "criminology": (113, 12, 1),
       "process management": (112, 12, 1), "library science": (103, 11, 1), "pattern recognition": (409, 44, 4),
       "statistical physics": (310, 32, 3), "seismology": (315, 34, 3), "marine engineering": (410, 45, 4),
       "cognitive psychology": (110, 12, 1), "ecology": (207, 23, 2), "nuclear engineering": (404, 42, 4),
       "economic policy": (112, 12, 1), "macroeconomics": (112, 12, 1), "reliability engineering": (402, 41, 4),
       "economic system": (112, 12, 1), "genetics": (201, 21, 2), "veterinary medicine": (207, 23, 2),
       "virology": (204, 22, 2), "petrology": (316, 34, 3), "business administration": (112, 12, 1),
       "computer science": (409, 44, 4), "medicinal chemistry": (205, 22, 2), "orthodontics": (205, 22, 2),
       "agricultural science": (207, 23, 2), "botany": (202, 21, 2), "medical education": (205, 22, 2),
       "physics": (307, 32, 3), "evolutionary biology": (203, 21, 2), "monetary economics": (112, 12, 1),
       "political economy": (112, 12, 1), "mathematical physics": (310, 32, 3), "quantum mechanics": (309, 32, 3),
       "positive economics": (112, 12, 1), "world wide web": (409, 44, 4), "mechanical engineering": (402, 41, 4),
       "economy": (112, 12, 1), "linguistics": (104, 11, 1), "social science": (111, 12, 1),
       "theoretical computer science": (409, 44, 4), "chromatography": (304, 31, 3), "geometry": (312, 33, 3),
       "construction engineering": (410, 45, 4), "nuclear physics": (309, 32, 3), "analytical chemistry": (304, 31, 3),
       "management science": (112, 12, 1), "astrophysics": (311, 32, 3), "forensic engineering": (406, 43, 4),
       "meteorology": (313, 34, 3), "physical geography": (317, 34, 3), "andrology": (205, 22, 2),
       "toxicology": (205, 22, 2), "ceramic materials": (405, 42, 4), "environmental science": (207, 23, 2),
       "computational chemistry": (304, 31, 3), "calculus": (312, 33, 3), "thermodynamics": (307, 32, 3),
       "immunology": (204, 22, 2), "demographic economics": (112, 12, 1), "economic history": (112, 12, 1),
       "crystallography": (316, 34, 3)}

boards = {
    101: 'Ancient Cultures',
    102: 'History',
    103: 'Fine Arts, Music, Theatre and Media Studies',
    104: 'Linguistics',
    105: 'Literary Studies',
    106: 'Social and Cultural Anthropology, Non-European Cultures, Jewish Studies and Religious Studies',
    107: 'Theology',
    108: 'Philosophy',
    109: 'Educational Research',
    110: 'Psychology',
    111: 'Social Sciences',
    112: 'Economics',
    113: 'Jurisprudence',
    201: 'Basic Biological and Medical Research',
    202: 'Plant Sciences',
    203: 'Zoology',
    204: 'Microbiology, Virology and Immunology',
    205: 'Medicine',
    206: 'Neurosciences',
    207: 'Agriculture, Forestry and Veterinary Medicine',
    301: 'Molecular Chemistry',
    302: 'Chemical Solid State and Surface Research',
    303: 'Physical and Theoretical Chemistry',
    304: 'Analytical Chemistry, Method Development (Chemistry)',
    305: 'Biological Chemistry and Food Chemistry',
    306: 'Polymer Research',
    307: 'Condensed Matter Physics',
    308: 'Optics, Quantum Optics and Physics of Atoms, Molecules and Plasmas',
    309: 'Particles, Nuclei and Fields',
    310: 'Statistical Physics, Soft Matter, Biological Physics, Nonlinear Dynamics',
    311: 'Astrophysics and Astronomy',
    312: 'Mathematics',
    313: 'Atmospheric Science, Oceanography and Climate Research',
    314: 'Geology and Palaeontology',
    315: 'Geophysics and Geodesy',
    316: 'Geochemistry, Mineralogy and Crystallography',
    317: 'Geography',
    318: 'Water Research',
    401: 'Production Technology',
    402: 'Mechanics and Constructive Mechanical Engineering',
    403: 'Process Engineering, Technical Chemistry',
    404: 'Heat Energy Technology, Thermal Machines, Fluid Mechanics',
    405: 'Materials Engineering',
    406: 'Materials Science',
    407: 'Systems Engineering',
    408: 'Electrical Engineering and Information Technology',
    409: 'Computer Science',
    410: 'Construction Engineering and Architecture'
}

areas = {
    11: 'Humanities',
    12: 'Social and Behavioural Sciences',
    21: 'Biology',
    22: 'Medicine',
    23: 'Agriculture, Forestry and Veterinary Medicine',
    31: 'Chemistry',
    32: 'Physics',
    33: 'Mathematics',
    34: 'Geosciences',
    41: 'Mechanical and Industrial Engineering',
    42: 'Thermal Engineering/Process Engineering',
    43: 'Materials Science and Engineering',
    44: 'Computer Science, Systems and Electrical Engineering',
    45: 'Construction Engineering and Architecture'
}

disciplines = {
    1: "Humanities and Social Sciences",
    2: "Life Sciences",
    3: "Natural Sciences",
    4: "Engineering Sciences"
}


def get_fos(x):
    if not x:
        return None
    fields = list({s.lower() for s in x}.intersection(set(FOS.keys())))
    if fields:
        x = list(zip(*[FOS[field] for field in fields]))
        return [
            list({boards.get(y) for y in x[0]}),
            list({areas.get(y) for y in x[1]}),
            list({disciplines.get(y) for y in x[2]})
        ]
    else:
        return None


def process(sc, input_path):
    df = (
        sc.read.json(input_path)
        .select(
            F.col("doi"),
            F.col("fos"),
            F.col("year")
        )
        .filter(~F.isnull("doi"))
        .filter(~F.isnull("year"))
        .withColumn(
            "fos_sampled",
            F.udf(
                get_fos,
                T.StructType([
                    T.StructField("board", T.ArrayType(T.StringType()),True),
                    T.StructField("area", T.ArrayType(T.StringType()),True),
                    T.StructField("discipline", T.ArrayType(T.StringType()),True)
                ])
            )("fos")
        )
        .select(
            F.col("doi"),
            F.col("year"),
            F.col("fos_sampled.*")
        )
        .coalesce(100)
    )
    return df


def run(sc, args):
    input_path = args[0]
    output_path = args[1]

    (
        process(sc, input_path)
        .write
        .mode("overwrite")
        .parquet(output_path)
    )


if __name__ == '__main__':
    INPUT_PATH = "file:/mnt/ceph/storage/corpora/corpora-thirdparty/corpus-microsoft-open-academic-graph-v1/*.txt"
    OUTPUT_PATH = "stereo-metadata.parquet"
    # spark session
    spark = (
        SparkSession
        .builder
        .config(conf=SparkConf())
        .getOrCreate()
    )
    # process
    (
        process(spark, INPUT_PATH)
        .write
        .mode("overwrite")
        .parquet(OUTPUT_PATH)
    )

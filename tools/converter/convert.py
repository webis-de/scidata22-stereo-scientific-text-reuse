import argparse
import lxml.etree as ET

def extract(inp: str):
    try:
        inp = bytes(inp, 'utf-8')
        tree = ET.fromstring(inp)
        for ref in tree.iter("{http://www.tei-c.org/ns/1.0}ref"):
            ref.text = ""
        text = []
        for elem in tree.iter("{http://www.tei-c.org/ns/1.0}p"):
            text.append(ET.tostring(elem, encoding="utf-8", method='text').decode("utf-8"))

        return " ".join(text)
    except Exception as e:
        print(e)
        return ""

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--input', type=str, required=True)
    parser.add_argument('--output', type=str, required=True)
    args = parser.parse_args()

    with open(args.input, "r") as infile:
        xml = "\n".join(infile.readlines())
        print(type(xml))
        text = extract(xml)

    with open(args.output, "w") as outfile:
        outfile.write(text)

if __name__ == "__main__":
    main()
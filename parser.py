import re



def rem_all_after(s, sep):
    return s.split(sep, 1)[0]


def rem_all_before(s, sep):
    index = s.index(sep)
    return s[index + 1:]


def sort_file(filename):
    fdr = open(filename, "r")
    lines = fdr.readlines()
    #print(lines)
    fdr.close()
    lines = list(set(lines))
    lines.sort()
    fdw = open(filename, "w")
    for line in lines:
        line = line.replace("ă","a");
        line = line.replace("î","i");
        line = line.replace("â","a");
        line = line.replace("ș","s");
        line = line.replace("ț","t");
        fdw.write(line)



def remove_short_words(s):
    b = ["Intranz", "Tranz", "pref", "tranz", "refl"]
    l = []
    for _s in s.split(" "):
        if len(_s) > 3 and _s not in b:
            l.append(_s)
    return " ".join(l)

#f = open("cuv5.txt", "r")
#out = open("out.txt", "w")

"""
line = f.readline()
cuv = []
while line != "":
    if "-" in line:
        line = line[:line.rindex("-")]
    line = line.split(";")

    for i in range(len(line)):
        line[i] = re.sub('\[.*?\]', '', line[i])
        line[i] = re.sub('\(.*?\)', '', line[i])
        line[i] = line[i].replace(".", " ")
        # line[i].find(", -")
        # cand ai cv -ceva vezi daca nu e alt cuvant da mi lene
        line[i] = ''.join([i.lower() for i in line[i] if (not i.isdigit() and i.isalpha() or i.isspace())])
        line[i] = remove_short_words(line[i])
        line[i] = line[i].replace("í", "i")
        line[i] = line[i].replace("á", "a")
        line[i] = line[i].replace("ó", "o")
        line[i] = line[i].replace("ú", "u")
        line[i] = line[i].replace("é", "e")
        line[i] = line[i].replace("ắ", "ă")
        line[i] = line[i].replace("ấ", "â")
        for l in line[i].split(" "):
            out.write(l)
            out.write("\n")
    line = f.readline()
f.close()
out.close()
"""
sort_file("ALL.txt")

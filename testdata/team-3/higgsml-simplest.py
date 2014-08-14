#!/usr/bin/env python

# completely self contained
# v2 20140522 : final ordering was buggy and too complicated. Simplified considerably

import csv
import optparse
import sys

# somewhat arbitrary value, should be optimised
threshold = -22.0

def printf(*args):
    if len(args) > 1: line = args[0]+"\n" % args[1:]
    else:             line = args[0]+"\n"
    return sys.stdout.write(line)

# compute AMS
def ams(s,b):
    from math import sqrt,log
    if b==0:
        return 0

    return sqrt(2*((s+b+10)*log(1+float(s)/(b+10))-s))

def main():
    rc = 1
    try:    rc = run()
    except: raise
    return rc

def run():
    global threshold

    parser = optparse.OptionParser(
        "This script computes the submission file with just a simple window on one variable DER_mass_MMC"
        )
    parser.add_option('--train', action="store_true", default=False, help="enable training mode")
    parser.add_option("--cut", default=threshold, help="cut-off value")

    opts, args = parser.parse_args()

    threshold = opts.cut
    
    if opts.train:
        return do_train(fname=args[0], trained=args[1])
    return run_prediction(fname=args[0], trained=args[1], ofname=args[2])


def do_train(fname="training.csv", trained="trained.dat"):
    printf("Reading in training file")
    alltraining = list(csv.reader(open(fname), delimiter=','))

    # first line is the list of variables
    headertraining        = alltraining[0]
    # cut off first line
    alltraining=alltraining[1:]

    # get the index of a few variables
    immc=headertraining.index("DER_mass_MMC")
    injet=headertraining.index("PRI_jet_num")
    iweight=headertraining.index("Weight")
    ilabel=headertraining.index("Label")
    iid=headertraining.index("EventId")
    
    
    printf("Loop on training dataset and compute the score")

    headertraining+=["myscore"]
    for entry in alltraining:
        # turn all entries from string to float, except EventId and PRI_jet_num to int, except label remains string
        for i in range(len(entry)):
            if not i in [ilabel,iid,injet]:
                entry[i]=float(entry[i])
            if i in [iid,injet]:
                entry[i]=int(entry[i])


        myscore=-abs(entry[immc]-125.) # this is a simple discriminating variable. Signal should be closer to zero.
        # minus sign so that signal has the highest values
        # so we will be making a simple window cut on the Higgs mass estimator
        # 125 GeV is the middle of the window
        entry+=[myscore]
        pass

    # at this stage alltraining is a list (one entry per line) of list of variables
    # which can be conveniently accessed by getting the index from the header 

    printf("Loop again to determine the AMS, using threshold:",threshold)
    sumsig=0.
    sumbkg=0.
    iscore=headertraining.index("myscore")
    for entry in alltraining:
        myscore=entry[iscore]
        entry+=[myscore]
        weight=entry[iweight]
        # sum event weight passing the selection. Of course in real life the threshold should be optimised
        if myscore >threshold:
            if entry[ilabel]=="s":
                sumsig+=weight
            else:
                sumbkg+=weight    
                pass
            pass
        pass
    
    # ok now we have our signal (sumsig) and background (sumbkg) estimation
    printf(" AMS computed from training file :",ams(sumsig,sumbkg),"( signal=",sumsig," bkg=",sumbkg,")")
    # delete big objects
    del alltraining
    

def run_prediction(fname="test.csv", trained="trained.dat", ofname="scores_test.csv"):
    printf("Reading in test file")
    alltest = list(csv.reader(open(fname), delimiter=','))
    headertest        = alltest[0]
    alltest=alltest[1:]


    printf("Compute the score for the test file entries ")

    # recompute variable indices for safety 
    immc=headertest.index("DER_mass_MMC")
    injet=headertest.index("PRI_jet_num")
    iid=headertest.index("EventId")
    headertest+=["myscore"]

    for entry in alltest:
        # turn all entries from string to float, except EventId and PRI_jet_num to int (there is no label)
        for i in range(len(entry)):
            if not i in [iid,injet]:
                entry[i]=float(entry[i])
            else:    
                entry[i]=int(entry[i])
        # add my score
        myscore=-abs(entry[immc]-125.)                                    
        entry+=[myscore]
        pass

    iscore=headertest.index("myscore")
    if iscore<0:
        printf("ERROR could not find variable myscore")
        raise Exception # should not happen

    printf("Sort on the score ") 
    # in the first version of the file, an auxilliary map was used, but this was useless
    alltestsorted=sorted(alltest,key=lambda entry: entry[iscore])
    # the RankOrder we want is now simply the entry number

    printf("Final loop to write the submission file %s", ofname)
    outputfile=open(ofname,"w")
    outputfile.write("EventId,RankOrder,Class\n")
    iid=headertest.index("EventId")
    if iid<0:
        printf("ERROR could not find variable EventId in test file")
        raise Exception # should not happen

    rank=1 # kaggle wants to start at 1
    for entry in alltestsorted:
        # compute label 
        slabel="b"
        if entry[iscore]>threshold: # arbitrary threshold
            slabel="s"

        outputfile.write(str(entry[iid])+",")
        outputfile.write(str(rank)+",")
        outputfile.write(slabel)            
        outputfile.write("\n")
        rank+=1
        pass

    outputfile.close()
    printf(" You can now submit %s to kaggle site", ofname)

    # delete big objects
    del alltest,alltestsorted
    return 0

if __name__ == "__main__":
    rc = main()
    sys.exit(rc)
    pass

